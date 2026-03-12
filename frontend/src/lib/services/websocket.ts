type EventCallback = (payload: any) => void;

interface SubscriptionMessage {
	action: 'subscribe' | 'unsubscribe';
	topic: string;
}

interface EventMessage {
	topic: string;
	payload: any;
}

export class WebSocketService {
	private ws: WebSocket | null = null;
	private readonly url: string;
	private subscribers: Map<string, Set<EventCallback>> = new Map();
	private reconnectAttempts = 0;
	private maxReconnectAttempts = 5;
	private reconnectDelay = 1000;
	private isIntentionallyClosed = false;

	constructor(url: string) {
		this.url = url;
	}

	connect(): Promise<void> {
		return new Promise((resolve, reject) => {
			if (this.ws?.readyState === WebSocket.OPEN) {
				resolve();
				return;
			}

			this.isIntentionallyClosed = false;
			this.ws = new WebSocket(this.url);

			this.ws.onopen = () => {
				console.log('WebSocket connected');
				this.reconnectAttempts = 0;

				// Resubscribe to all topics
				this.subscribers.forEach((_, topic) => {
					this.sendSubscription('subscribe', topic);
				});

				resolve();
			};

			this.ws.onmessage = (event) => {
				try {
					const message: EventMessage = JSON.parse(event.data);
					this.handleMessage(message);
				} catch (err) {
					console.error('Failed to parse WebSocket message:', err);
				}
			};

			this.ws.onerror = (error) => {
				console.error('WebSocket error:', error);
				reject(error);
			};

			this.ws.onclose = () => {
				console.log('WebSocket closed');
				this.ws = null;

				if (!this.isIntentionallyClosed) {
					this.attemptReconnect();
				}
			};
		});
	}

	subscribe(topic: string, callback: EventCallback): () => void {
		if (!this.subscribers.has(topic)) {
			this.subscribers.set(topic, new Set());

			// Send subscription message if connected
			if (this.ws?.readyState === WebSocket.OPEN) {
				this.sendSubscription('subscribe', topic);
			}
		}

		this.subscribers.get(topic)!.add(callback);

		// Return unsubscribe function
		return () => this.unsubscribe(topic, callback);
	}

	disconnect() {
		this.isIntentionallyClosed = true;
		this.subscribers.clear();

		if (this.ws) {
			this.ws.close();
			this.ws = null;
		}
	}

	isConnected(): boolean {
		return this.ws?.readyState === WebSocket.OPEN;
	}

	private attemptReconnect() {
		if (this.reconnectAttempts < this.maxReconnectAttempts) {
			this.reconnectAttempts++;
			const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);

			console.log(
				`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts}) in ${delay}ms`
			);

			setTimeout(() => {
				this.connect().catch((err) => {
					console.error('Reconnection failed:', err);
				});
			}, delay);
		} else {
			console.error('Max reconnection attempts reached');
		}
	}

	private unsubscribe(topic: string, callback: EventCallback) {
		const callbacks = this.subscribers.get(topic);
		if (!callbacks) return;

		callbacks.delete(callback);

		if (callbacks.size === 0) {
			this.subscribers.delete(topic);

			// Send unsubscription message if connected
			if (this.ws?.readyState === WebSocket.OPEN) {
				this.sendSubscription('unsubscribe', topic);
			}
		}
	}

	private sendSubscription(action: 'subscribe' | 'unsubscribe', topic: string) {
		const message: SubscriptionMessage = { action, topic };
		this.send(message);
	}

	private send(data: any) {
		if (this.ws?.readyState === WebSocket.OPEN) {
			this.ws.send(JSON.stringify(data));
		} else {
			console.warn('WebSocket is not open. Message not sent:', data);
		}
	}

	private handleMessage(message: EventMessage) {
		const callbacks = this.subscribers.get(message.topic);
		if (callbacks) {
			callbacks.forEach((callback) => {
				try {
					callback(message.payload);
				} catch (err) {
					console.error(`Error in callback for topic ${message.topic}:`, err);
				}
			});
		}
	}
}

// Singleton instance for the application
let wsInstance: WebSocketService | null = null;

export function getWebSocketService(): WebSocketService {
	if (!wsInstance) {
		wsInstance = new WebSocketService('ws://localhost:3000/api/events/ws');
	}
	return wsInstance;
}
