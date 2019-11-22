import React from 'react';

const FetchingComponent = ({ isLoading, isError, render }) => {
    return isLoading || isError ? (
        <div>
            {isLoading
                ? 'Loading...'
                : 'Error occured while loading data. Are you sure server is running?'}
        </div>
    ) : (
        render()
    );
};

export default FetchingComponent;
