export interface Qualifying {
    status: string;
    players: QualifyingPlayer[];
}

export interface QualifyingPlayer {
    name: string;
    bestTime: string;
    position: number;
    signupDate: string;
}