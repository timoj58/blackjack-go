import { Link } from 'react-router-dom';

function TableTile(props) {
  const {
    tableDetails: { table, cut, stake, players },
    playerId
  } = props;
  return (
    <div style={{ borderStyle: 'solid', borderRadius: '10px', padding: '10px', width: '800px' }}>
      <div className="row">
        <h1>{table}</h1>
        <div className="column">
          <p style={{ textAlign: 'right' }}>Stake: £{stake}</p>
          <p style={{ textAlign: 'right' }}>Cut: {cut}</p>
          <p style={{ textAlign: 'right' }}>Players: {players}</p>
          <Link to="/table" state={{ playerId, tableId: table }}>
            <p style={{ textAlign: 'right' }}>Join</p>
          </Link>
        </div>
      </div>
    </div>
  );
}

export default TableTile;
