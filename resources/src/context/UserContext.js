import React from 'react';

const notImplemented = () => new Error('not implemented yet');

export default React.createContext({
    isAdmin:notImplemented,
})
