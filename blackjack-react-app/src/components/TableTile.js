import { Link } from 'react-router-dom';

function TableTile(props) {
  const {
    tableDetails: { id, name, cut, stake, players, status },
    playerId
  } = props;
  return (
    <div style={{ borderStyle: 'solid', borderRadius: '10px', padding: '10px', width: '33vh' }}>
      <div className="row">
        <h1>Table {name}</h1>
        <div className="column">
          <p style={{ textAlign: 'right' }}>
            Stake - Cut: £{stake} - {cut}
          </p>
          <p style={{ textAlign: 'right' }}>Players: {players}</p>
          {status === true ? (
            <Link to="/table" state={{ playerId, tableId: id }}>
              <p style={{ textAlign: 'right' }}>Join</p>
            </Link>
          ) : (
            <p style={{ textAlign: 'right' }}>In play</p>
          )}
        </div>
      </div>
    </div>
  );
}

export default TableTile;
