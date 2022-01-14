import { useLocation } from 'react-router-dom';
import React, { useState, useEffect } from 'react';
import { get } from '../util/Websocket';
import TableTile from '../components/TableTile';

const socket = get();

function Tables() {
  const location = useLocation();
  const [tables, setTables] = useState([]);
  const [playerId, setPlayerId] = useState(location.state.playerId);

  useEffect(() => {
    if (socket.readyState === WebSocket.OPEN) {
      socket.send(
        JSON.stringify({
          playerId,
          action: 'list',
          data: ''
        })
      );
    }
  }, []);

  useEffect(() => {
    socket.onmessage = (event) => {
      const data = event.data.split('\n').map((d) => JSON.parse(d));
      let addTables = false;
      let toAdd = [];
      data.forEach((element) => {
        if (element.type === 'player') {
          setPlayerId(JSON.parse(event.data).playerId);
          socket.send(
            JSON.stringify({
              playerId: JSON.parse(event.data).playerId,
              action: 'list',
              data: ''
            })
          );
        } else if (element.type === 'tables') {
          addTables = true;
          toAdd = toAdd.concat(element);
        }
      });
      if (addTables) {
        setTables(tables.concat(toAdd));
      }
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
