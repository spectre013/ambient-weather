export interface Alerts {
    alerts: Alert[]
}
export interface Alert {
    id: number;
    alertid: string;
    wxtype: string;
    areadesc: string;
    sent: string;
    effective: string;
    onset: string;
    expires: string;
    end: string;
    status: string;
    messagetype: string;
    category: string;
    severity: string;
    certainty: string;
    urgency: string;
    event: string;
    sender: string;
    senderName: string;
    headline: string;
    description: string;
    instruction: string;
    response: string;
}