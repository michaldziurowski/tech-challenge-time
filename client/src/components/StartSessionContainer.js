import React, { useState } from 'react';

const StartSessionContainer = ({ onStartSession }) => {
    const [name, setName] = useState('');

    return (
        <div>
            <input value={name} onChange={e => setName(e.target.value)} />
            <button onClick={() => onStartSession(name)}>Start</button>
        </div>
    );
};

export default StartSessionContainer;
