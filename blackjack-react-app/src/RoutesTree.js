import React from 'react';
import { Route, Routes } from 'react-router-dom';
import Home from './screens/Home';
import Tables from './screens/Tables';
import Table from './screens/Table';

function RoutesTree() {
  return (
    <div>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/tables" element={<Tables />} />
        <Route path="/table" element={<Table />} />
      </Routes>
    </div>
  );
}

export default RoutesTree;
