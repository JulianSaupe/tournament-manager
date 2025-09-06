export interface Qualifying {
    status: string;
    players: QualifyingPlayer[];
}

export interface QualifyingPlayer {
    name: string;
    time: number;
    position: number;
    signupDate: string;
}

export enum QualifyingFilter {
    ALL = 'all',
    QUALIFIED = 'qualified',
    UNQUALIFIED = 'unqualified',
}