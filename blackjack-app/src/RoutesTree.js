import React from 'react';
import { Route, Routes } from 'react-router-dom';
import Home from './screens/Home';

function RoutesTree() {
  return (
    <div>
      <Routes>
        <Route path="/" element={<Home />} />
      </Routes>
    </div>
  );
}

export default RoutesTree;
