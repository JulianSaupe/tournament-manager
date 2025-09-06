import type { PageServerLoad } from './$types';

export type IndexTournament = {
  id: string;
  name: string;
  description: string;
  startDate: string;
  endDate: string;
  status: string; // Backend returns uppercase: DRAFT | ACTIVE | COMPLETED | CANCELLED
  playerCount: number;
};

type StandardResponse<T> = {
  server: string;
  startTime: string;
  endTime: string;
  statusCode: number;
  data?: T;
};

export const load: PageServerLoad = async ({ fetch }) => {
  try {
    const res = await fetch('http://localhost:3000/api/tournament');
    if (!res.ok) {
      const text = await res.text().catch(() => res.statusText);
      return { tournaments: [] as IndexTournament[], error: `Backend error: ${text || res.status}` };
    }
    const body = (await res.json()) as StandardResponse<unknown>;
    const raw = (body as StandardResponse<unknown>).data;
    const tournaments = Array.isArray(raw) ? (raw as IndexTournament[]) : [];
    return { tournaments, error: null as string | null };
  } catch (err) {
    return {
      tournaments: [] as IndexTournament[],
      error: err instanceof Error ? err.message : 'Unknown error'
    };
  }
};
