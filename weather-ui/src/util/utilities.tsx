import moment from "moment";

// @ts-expect-error - ignore function parameter types
export const CustomTooltip = ({ active, payload }) => {
    if (active && payload && payload.length) {
        const date = new Date(payload[0].payload.date);
        const formattedDate = date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
        });

        return (
            <div style={{
                backgroundColor: '#374151',
                border: 'none',
                borderRadius: '0.5rem',
                padding: '1rem',
                color: '#e5e7eb'
            }}>
                <p className="label" style={{ fontWeight: 'bold' }}>{formattedDate}</p>
                <p className="intro" style={{ color: payload[0].color }}>{`${payload[0].name} : ${payload[0].value} in`}</p>
            </div>
        );
    }
    return null;
};

// @ts-expect-error - ignore function parameter types
export const CustomForecastTooltip = ({ active, payload }) => {
    if (active && payload && payload.length) {
        const date = moment.unix(payload[0].payload.datetimeEpoch);
        const formattedDate = date.format('HH:mm:ss');

        return (
            <div style={{
                backgroundColor: '#374151',
                border: 'none',
                borderRadius: '0.5rem',
                padding: '1rem',
                color: '#e5e7eb'
            }}>
                <p className="label" style={{ fontWeight: 'bold' }}>{formattedDate}</p>
                <p className="intro" style={{ color: payload[0].color }}>{`${payload[0].name} : ${payload[0].value} in`}</p>
            </div>
        );
    }
    return null;
};