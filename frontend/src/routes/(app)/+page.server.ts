import type {PageServerLoad} from './$types';
import {TournamentStatus} from "$lib/types/tournament/tournament";

export type IndexTournament = {
    id: string;
    name: string;
    description: string;
    startDate: string;
    endDate: string;
    status: TournamentStatus;
    playerCount: number;
};

type StandardResponse<T> = {
    server: string;
    startTime: string;
    endTime: string;
    statusCode: number;
    data?: T;
};

export const load: PageServerLoad = async ({fetch}) => {
    try {
        const response = await fetch('http://localhost:3000/api/tournament');

        if (!response.ok) {
            const text = await response.text().catch(() => response.statusText);
            return {
                tournaments: [] as IndexTournament[],
                error: `Backend error: ${text || response.status}`
            };
        }

        const body = (await response.json()) as StandardResponse<unknown>;
        const raw = (body as StandardResponse<unknown>).data;
        const tournaments = Array.isArray(raw) ? raw.map((tournament: any): IndexTournament => {
            const statusValue = tournament.status?.toLowerCase();
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

            return {
                id: tournament.id,
                name: tournament.name,
                description: tournament.description,
                startDate: tournament.startDate,
                endDate: tournament.endDate,
                status: parsedStatus,
                playerCount: tournament.playerCount || 0
            };
        }) : [];


        return {
            tournaments: tournaments,
            error: null as string | null
        };
    } catch (err) {
        return {
            tournaments: [] as IndexTournament[],
            error: err instanceof Error ? err.message : 'Unknown error'
        };
    }
};
