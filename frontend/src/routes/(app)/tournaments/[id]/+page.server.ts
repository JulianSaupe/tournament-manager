import type {PageServerLoad} from "./$types";
import type {Tournament} from "$lib/types/tournament/tournament";
import {TournamentStatus} from "$lib/types/tournament/tournament";
import {type Qualifying, type QualifyingPlayer} from "$lib/types/tournament/qualifying";

type StandardResponse<T> = {
    server: string;
    startTime: string;
    endTime: string;
    statusCode: number;
    data?: T;
};

export const load: PageServerLoad = async ({params}) => {
    try {
        const tournamentResponse = await fetch(`http://localhost:3000/api/tournament/${params.id}`);

        if (!tournamentResponse.ok) {
            const text = await tournamentResponse.text().catch(() => tournamentResponse.statusText);
            return {
                error: `Backend error: ${text || tournamentResponse.status}`
            };
        }

        let body = (await tournamentResponse.json()) as StandardResponse<unknown>;
        const rawTournament = (body as StandardResponse<Tournament>).data;

        if (rawTournament === undefined) {
            return {
                error: "Couldn't fetch tournament",
            }
        }

        const statusValue = rawTournament.status?.toLowerCase();
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
            id: rawTournament.id,
            name: rawTournament.name,
            description: rawTournament.description,
            startDate: rawTournament.startDate,
            endDate: rawTournament.endDate,
            status: parsedStatus,
            playerCount: rawTournament.playerCount || 0
        };

        const qualifyingResponse = await fetch(`http://localhost:3000/api/tournament/${params.id}/qualifying`);

        body = (await qualifyingResponse.json()) as StandardResponse<unknown>;
        const rawQualifying = (body as StandardResponse<Qualifying>).data;

        if (rawQualifying === undefined) {
            return {
                error: "Couldn't fetch qualifying",
            }
        }

        let players: QualifyingPlayer[] = rawQualifying.players.map((player: any): QualifyingPlayer => {
            return {
                name: player['name'],
                position: player['position'],
                signupDate: player['signup_date'],
                time: player['time'],
            }
        })

        let qualifying: Qualifying = {
            status: 'test',
            players: players,
        };

        return {
            tournament: tournament,
            qualifying: qualifying,
        };
    } catch (err) {
        return {
            error: err instanceof Error ? err.message : 'Unknown error'
        };
    }
};