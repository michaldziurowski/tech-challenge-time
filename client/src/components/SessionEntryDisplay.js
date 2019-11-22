import React from 'react';

const displayDurationFromNanoseconds = nanosecondDuration => {
    let initialSeconds = Math.floor(nanosecondDuration / 1000000000);
    let hours = Math.floor(initialSeconds / 3600);
    let minutes = Math.floor((initialSeconds - hours * 3600) / 60);
    let seconds = initialSeconds - hours * 3600 - minutes * 60;

    return `${hours}h${minutes}m${seconds}s`;
};

const SessionEntryDisplay = ({ session, duration, onToggleSession }) => {
    return (
        <div key={session.SessionId}>
            {session.Name} {displayDurationFromNanoseconds(duration)}{' '}
            {session.IsOpen ? (
                <button
                    onClick={() => onToggleSession(session.SessionId, false)}
                >
                    Close
                </button>
            ) : (
                <button
                    onClick={() => onToggleSession(session.SessionId, true)}
                >
                    Resume
                </button>
            )}
        </div>
    );
};

export default SessionEntryDisplay;
