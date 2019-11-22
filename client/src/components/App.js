import React, { useState, useEffect } from 'react';
import './App.css';
import FetchingComponent from './FetchingComponent';
import {
    fetchSessions,
    startSessionRequest,
    stopSessionRequest,
    resumeSessionRequest
} from '../api';
import StartSessionContainer from './StartSessionContainer';
import { DAY_MODE, WEEK_MODE, MONTH_MODE } from '../consts';
import SessionsHistory from './SessionsHistory';

const today = () => new Date();
const todayStart = new Date(today().setHours(0, 0, 0, 0));
const todayEnd = new Date(today().setHours(23, 59, 59, 999));
const lastWeekStart = new Date(
    new Date(today().setDate(today().getDate() - 7)).setHours(0, 0, 0, 0)
);
const lastMonthStart = new Date(
    new Date(today().setMonth(today().getMonth() - 1)).setHours(0, 0, 0, 0)
);

const getDateRangeForMode = mode => {
    switch (mode) {
        case DAY_MODE:
            return [todayStart, todayEnd];
        case WEEK_MODE:
            return [lastWeekStart, todayEnd];
        case MONTH_MODE:
            return [lastMonthStart, todayEnd];
        default:
            return [];
    }
};

const App = () => {
    const [mode, setMode] = useState(DAY_MODE);
    const [from, to] = getDateRangeForMode(mode);
    const [sessions, setSessions] = useState();
    const [isLoading, setIsLoading] = useState(false);
    const [isError, setIsError] = useState(false);

    async function getData() {
        try {
            setIsLoading(true);
            const data = await fetchSessions(from, to);
            setIsLoading(false);
            setSessions(data);
        } catch {
            setIsLoading(false);
            setIsError(true);
        }
    }

    async function startSession(name) {
        try {
            setIsLoading(true);
            await startSessionRequest(name, new Date());
            setIsLoading(false);
            getData();
        } catch {
            setIsLoading(false);
            setIsError(true);
        }
    }

    async function toggleSession(sessionId, shouldOpen) {
        try {
            setIsLoading(true);
            if (shouldOpen) {
                await resumeSessionRequest(sessionId, new Date());
            } else {
                await stopSessionRequest(sessionId, new Date());
            }
            setIsLoading(false);
            await getData();
        } catch {
            setIsLoading(false);
            setIsError(true);
        }
    }

    useEffect(() => {
        getData();
    }, [mode]);

    return (
        <div className="App">
            <header className="App-header">
                <p>Time tracking</p>
                <FetchingComponent
                    isLoading={isLoading}
                    isError={isError}
                    render={() => (
                        <div>
                            <StartSessionContainer
                                onStartSession={name => startSession(name)}
                            />
                            <SessionsHistory
                                sessions={sessions}
                                onSetMode={setMode}
                                onToggleSession={toggleSession}
                            />
                        </div>
                    )}
                />
            </header>
        </div>
    );
};

export default App;
