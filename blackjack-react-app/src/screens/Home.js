import React from 'react';
import { Link } from 'react-router-dom';
import { init } from '../util/Websocket';

const socket = init();

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      playerId: ''
    };
  }

  componentDidMount() {
    socket.onmessage = (event) => {
      this.setState({ playerId: JSON.parse(event.data).playerId });
    };
  }

  render() {
    const { playerId } = this.state;
    return (
      <form>
        <Link to="/tables" state={{ playerId }}>
          Enter
        </Link>
      </form>
    );
  }
}

export default Home;
