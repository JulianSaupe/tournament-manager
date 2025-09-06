import type {PageServerLoad} from "./$types";
import type {Tournament} from "$lib/types/tournament/tournament";
import {TournamentStatus} from "$lib/types/tournament/tournament";

type StandardResponse<T> = {
    server: string;
    startTime: string;
    endTime: string;
    statusCode: number;
    data?: T;
};

export const load: PageServerLoad = async ({params}) => {
    try {
        const response = await fetch('http://localhost:3000/api/tournament/' + params.id);

        if (!response.ok) {
            const text = await response.text().catch(() => response.statusText);
            return {
                error: `Backend error: ${text || response.status}`
            };
        }

        const body = (await response.json()) as StandardResponse<unknown>;
        const raw = (body as StandardResponse<Tournament>).data;

        if (raw === undefined) {
            return {
                error: "Couldn't fetch tournament",
            }
        }

        const statusValue = raw.status?.toLowerCase();
        let parsedStatus: TournamentStatus;

        switch (statusValue) {
            case 'draft':
                parsedStatus = TournamentStatus.DRAFT;
                break;
            case 'active':
                parsedStatus = TournamentStatus.ACTIVE;
                break;
            case 'completed':
                parsedStatus = TournamentStatus.COMPLETED;
                break;
            case 'cancelled':
                parsedStatus = TournamentStatus.CANCELLED;
                break;
            default:
                parsedStatus = TournamentStatus.DRAFT;
                break;
        }

        let tournament: Tournament = {
            id: raw.id,
            name: raw.name,
            description: raw.description,
            startDate: raw.startDate,
            endDate: raw.endDate,
            status: parsedStatus,
            playerCount: raw.playerCount || 0
        };


        return {
            tournament: tournament,
        };
    } catch (err) {
        return {
            error: err instanceof Error ? err.message : 'Unknown error'
        };
    }
};