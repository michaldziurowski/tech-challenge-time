const serverUrl = 'http://localhost:8080/api/v1';

export const fetchSessions = async (startDate, endDate) => {
    const response = await fetch(
        `${serverUrl}/sessions?from=${startDate.toISOString()}&to=${endDate.toISOString()}`
    );
    return await response.json();
};

export const startSessionRequest = (name, startTime) => {
    return fetch(`${serverUrl}/sessions`, {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name: name, time: startTime.toISOString() })
    });
};

export const stopSessionRequest = (sessionId, stopTime) => {
    return fetch(`${serverUrl}/sessions/${sessionId}/stop`, {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ time: stopTime.toISOString() })
    });
};

export const resumeSessionRequest = (sessionId, resumeTime) => {
    return fetch(`${serverUrl}/sessions/${sessionId}/resume`, {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ time: resumeTime.toISOString() })
    });
};
