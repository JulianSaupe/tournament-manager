import {TournamentStatus} from "$lib/types/tournament/tournament";

type StatusConfig = {
    name: string;
    color: string;
    bgColor: string;
    icon: string;
};

export const statusConfig: Record<TournamentStatus, StatusConfig> = {
    draft: {
        name: 'Draft',
        color: 'badge-neutral',
        bgColor: 'bg-neutral/10',
        icon:
            '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />'
    },
    active: {
        name: 'Active',
        color: 'badge-success',
        bgColor: 'bg-success/10',
        icon:
            '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />'
    },
    completed: {
        name: 'Completed',
        color: 'badge-info',
        bgColor: 'bg-info/10',
        icon:
            '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />'
    },
    cancelled: {
        name: 'Cancelled',
        color: 'badge-error',
        bgColor: 'bg-error/10',
        icon:
            '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />'
    }
};