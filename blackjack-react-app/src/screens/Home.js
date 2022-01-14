import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { init } from '../util/Websocket';

const socket = init();

function Home() {
  const [playerId, setPlayerId] = useState('');

  useEffect(() => {
    socket.onmessage = (event) => {
      const data = event.data.split('\n').map((d) => JSON.parse(d));

      data.forEach((element) => {
        if (element.type === 'player') {
          setPlayerId(JSON.parse(event.data).playerId);
        }
      });
    };
  }, []);

  return (
    <form>
      <Link to="/tables" state={{ playerId }}>
        Enter
      </Link>
    </form>
  );
}

export default Home;
