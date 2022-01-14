import { useLocation, useNavigate } from 'react-router-dom';
import React, { useState, useCallback, useEffect } from 'react';
import { get } from '../util/Websocket';
import Hand from '../components/Hand';

const socket = get();

function Table() {
  const location = useLocation();
  const navigate = useNavigate();

  const [playerId] = useState(location.state.playerId);
  const [tableId] = useState(location.state.tableId);

  const [dealerCards, setDealerCards] = useState([]);
  const [playerCards, setPlayerCards] = useState([]);

  const [message, setMessage] = useState('');
  const [cards, setCards] = useState({});
  const [players, setPlayers] = useState([]);

  const hit = useCallback(() => {
    socket.send(
      JSON.stringify({
        playerId,
        action: 'hit',
        data: tableId
      })
    );
  }, []);

  const stick = useCallback(() => {
    socket.send(
      JSON.stringify({
        playerId,
        action: 'stick',
        data: tableId
      })
    );
  }, []);

  const leave = () => {
    socket.send(
      JSON.stringify({
        playerId,
        action: 'leave',
        data: tableId
      })
    );
    navigate('/tables', { state: { playerId } });
  };

  useEffect(() => {
    if (socket.readyState === WebSocket.OPEN) {
      socket.send(
        JSON.stringify({
          playerId,
          action: 'join',
          data: tableId
        })
      );
    } else {
      navigate('/');
    }
  }, []);

  useEffect(() => {
    window.onbeforeunload = function () {
      return true;
    };

    return () => {
      window.onbeforeunload = null;
    };
  }, []);

  useEffect(() => {
    socket.onmessage = (event) => {
      const data = event.data.split('\n').map((d) => JSON.parse(d));

      data.forEach((element) => {
        if (element.type === 'result') {
          setMessage(element.data);
          setPlayerCards([]);
          setDealerCards([]);
          players.forEach((p) => {
            cards[p.id] = [];
            setCards(cards);
          });
        } else if (element.type === 'message' || element.type === 'game') {
          setMessage(element.data);
        } else if (element.type === 'card') {
          if (element.dealer === true) {
            setDealerCards(dealerCards.concat(element));
          } else if (element.id === playerId) {
            setPlayerCards(playerCards.concat(element));
          } else {
            if (cards[element.id] === undefined) {
              cards[element.id] = [];
              setPlayers(players.concat({ id: element.id, dealer: element.dealer }));
            }
            cards[element.id] = cards[element.id].concat(element);
            setCards(cards);
          }
        }
      });
    };
  });

  return (
    <form>
      <div className="row">
        <p>{message}</p>
        <ul style={{ listStyleType: 'none' }}>
          <li style={{ marginBottom: '10px' }} key="dealer">
            <Hand title="dealer" cards={dealerCards} />
          </li>
          <li style={{ marginBottom: '10px' }} key={playerId}>
            <Hand title="you" cards={playerCards} />
          </li>
          <div style={{ marginTop: '10px', marginBottom: '10px' }}>
            <button type="button" onClick={hit}>
              Hit
            </button>
            <button type="button" onClick={stick}>
              Stick
            </button>
            <button type="button" onClick={leave}>
              Leave
            </button>
          </div>
          {players.map((p) => (
            <li style={{ marginBottom: '10px' }} key={p.id}>
              <Hand title={p.id} cards={cards[p.id]} />
            </li>
          ))}
        </ul>
      </div>
    </form>
  );
}

export default Table;
