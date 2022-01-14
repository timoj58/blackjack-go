import { useLocation } from 'react-router-dom';
import React, { useState, useEffect } from 'react';
import { get } from '../util/Websocket';
import TableTile from '../components/TableTile';

const socket = get();

function Tables() {
  const location = useLocation();
  const [tables, setTables] = useState([]);
  const [playerId] = useState(location.state.playerId);

  useEffect(() => {
    socket.send(
      JSON.stringify({
        playerId,
        action: 'list',
        data: ''
      })
    );
  }, []);

  useEffect(() => {
    socket.onmessage = (event) => {
      const data = event.data.split('\n').map((d) => JSON.parse(d));
      setTables(tables.concat(data));
    };
  });

  return (
    <form>
      <ul style={{ listStyleType: 'none' }}>
        {tables.map((d) => (
          <li style={{ marginBottom: '10px' }} key={d.table}>
            <TableTile tableDetails={d} playerId={playerId} />
          </li>
        ))}
      </ul>
    </form>
  );
}

export default Tables;
