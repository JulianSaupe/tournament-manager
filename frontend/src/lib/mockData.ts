// Tournament data structure
export interface Tournament {
  id: string;
  name: string;
  startDate: string; // ISO date string
  endDate: string;   // ISO date string
  playerCount: number;
  status: 'draft' | 'active' | 'completed' | 'cancelled';
}

// Mock tournament data
export const mockTournaments: Tournament[] = [
  // Draft tournaments
  {
    id: '1',
    name: 'Summer Chess Championship',
    startDate: '2025-09-01T00:00:00Z',
    endDate: '2025-09-15T00:00:00Z',
    playerCount: 16,
    status: 'draft'
  },
  {
    id: '2',
    name: 'Regional Tennis Open',
    startDate: '2025-10-05T00:00:00Z',
    endDate: '2025-10-12T00:00:00Z',
    playerCount: 32,
    status: 'draft'
  },
  {
    id: '3',
    name: 'Local Poker Tournament',
    startDate: '2025-08-20T00:00:00Z',
    endDate: '2025-08-21T00:00:00Z',
    playerCount: 8,
    status: 'draft'
  },
  
  // Active tournaments
  {
    id: '4',
    name: 'City Basketball League',
    startDate: '2025-07-15T00:00:00Z',
    endDate: '2025-08-30T00:00:00Z',
    playerCount: 12,
    status: 'active'
  },
  {
    id: '5',
    name: 'Online Gaming Championship',
    startDate: '2025-08-01T00:00:00Z',
    endDate: '2025-08-15T00:00:00Z',
    playerCount: 64,
    status: 'active'
  },
  {
    id: '6',
    name: 'Corporate Table Tennis',
    startDate: '2025-07-25T00:00:00Z',
    endDate: '2025-08-10T00:00:00Z',
    playerCount: 24,
    status: 'active'
  },
  
  // Completed tournaments
  {
    id: '7',
    name: 'Spring Soccer Cup',
    startDate: '2025-04-10T00:00:00Z',
    endDate: '2025-05-15T00:00:00Z',
    playerCount: 16,
    status: 'completed'
  },
  {
    id: '8',
    name: 'Winter Chess Tournament',
    startDate: '2025-01-05T00:00:00Z',
    endDate: '2025-01-20T00:00:00Z',
    playerCount: 8,
    status: 'completed'
  },
  {
    id: '9',
    name: 'Annual Golf Classic',
    startDate: '2025-06-01T00:00:00Z',
    endDate: '2025-06-05T00:00:00Z',
    playerCount: 32,
    status: 'completed'
  },
  
  // Cancelled tournaments
  {
    id: '10',
    name: 'Beach Volleyball Series',
    startDate: '2025-07-10T00:00:00Z',
    endDate: '2025-07-12T00:00:00Z',
    playerCount: 16,
    status: 'cancelled'
  },
  {
    id: '11',
    name: 'Charity Marathon Race',
    startDate: '2025-05-20T00:00:00Z',
    endDate: '2025-05-20T00:00:00Z',
    playerCount: 100,
    status: 'cancelled'
  }
];

// Helper function to group tournaments by status
export function getTournamentsByStatus(tournaments: Tournament[]): Record<string, Tournament[]> {
  return tournaments.reduce((acc, tournament) => {
    if (!acc[tournament.status]) {
      acc[tournament.status] = [];
    }
    acc[tournament.status].push(tournament);
    return acc;
  }, {} as Record<string, Tournament[]>);
}

// Helper function to format dates
export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  });
}