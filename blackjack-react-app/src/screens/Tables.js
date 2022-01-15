import { useLocation } from 'react-router-dom';
import React, { useState, useEffect } from 'react';
import { get } from '../util/Websocket';
import TableTile from '../components/TableTile';

const socket = get();

function Tables() {
  const location = useLocation();
  const [tables, setTables] = useState({});
  const [tablesList, setTablesList] = useState([]);
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
          tables[element.id] = element;
          setTables(tables);
          setTablesList(Object.values(tables));
        }
      });
    };
  });

  return (
    <form>
      <div className="grid">
        {tablesList.map((d) => (
          <div style={{ margin: '10px' }} key={d.id}>
            <TableTile tableDetails={d} playerId={playerId} />
          </div>
        ))}
      </div>
    </form>
  );
}

export default Tables;
