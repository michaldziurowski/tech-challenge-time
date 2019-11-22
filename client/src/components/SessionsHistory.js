import React from 'react';
import { DAY_MODE, WEEK_MODE, MONTH_MODE } from '../consts';
import SessionEntryDisplayContainer from './SessionEntryDisplayContainer';

const SessionsHistory = ({ sessions, onSetMode, onToggleSession }) => (
    <div>
        <div>
            <button onClick={() => onSetMode(DAY_MODE)}>Day</button>
            <button onClick={() => onSetMode(WEEK_MODE)}>Week</button>
            <button onClick={() => onSetMode(MONTH_MODE)}>Month</button>
        </div>
        <div>
            {(sessions &&
                sessions.map(s => (
                    <SessionEntryDisplayContainer
                        session={s}
                        onToggleSession={onToggleSession}
                        key={s.SessionId}
                    />
                ))) ||
                null}
        </div>
    </div>
);

export default SessionsHistory;
