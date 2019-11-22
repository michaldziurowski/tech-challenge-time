import React, { useState, useEffect } from 'react';
import SessionEntryDisplay from './SessionEntryDisplay';

const SECOND_IN_MILLISECONDS = 1000;
const SECOND_IN_NANOSECONDS = 1000000000;

const SessionEntryDisplayContainer = ({ session, onToggleSession }) => {
    const [duration, setDuration] = useState(session.Duration);

    useEffect(() => {
        let intervalId;
        if (session.IsOpen) {
            intervalId = setInterval(() => {
                // naive assumption that interval is run every second but for the sake of demo should be enough
                setDuration(prev => prev + SECOND_IN_NANOSECONDS);
            }, SECOND_IN_MILLISECONDS);
        } else {
            setDuration(session.Duration);
        }

        return () => {
            if (intervalId) {
                clearInterval(intervalId);
            }
        };
    }, [session, duration]);

    return (
        <SessionEntryDisplay
            session={session}
            duration={duration}
            onToggleSession={onToggleSession}
        />
    );
};

export default SessionEntryDisplayContainer;
