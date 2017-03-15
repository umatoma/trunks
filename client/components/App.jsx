import React from 'react';
import { OrderedSet, Record } from 'immutable';
import Header from './Header';
import FromPostAttack from './FormPostAttack';
import Metrics from './Metrics';
import Notifications from './Notifications';

const MetricsModel = Record({
  is_active: false,
  bytes_in: {},
  bytes_out: {},
  duration: 0,
  earliest: '',
  end: '',
  errors: [],
  latencies: {},
  latest: {},
  rate: 0,
  requests: 0,
  status_codes: {},
  success: 0,
  wait: 0,
});

class App extends React.Component {
  constructor() {
    super();

    // create websocket connection
    const conn = new WebSocket(`ws://${document.location.host}/ws`);
    conn.onclose = () => this.handleCloseWebSocket();
    conn.onmessage = evt => this.handleWebSocketMessage(evt);

    this.state = {
      webSocketConn: conn,
      notifications: OrderedSet(),
      metrics: new MetricsModel(),
    };

    this.addNotify = this.addNotify.bind(this);
    this.handleCloseNotify = this.handleCloseNotify.bind(this);
  }

  handleCloseWebSocket() {
    this.addNotify('WebSocket connection closed');
  }

  handleWebSocketMessage(evt) {
    const { event, data } = JSON.parse(evt.data);
    console.log(event, data);
    switch (event) {
      case 'attackMetrics':
        this.setState({
          metrics: new MetricsModel(Object.assign({ is_active: true }, data)),
        });
        break;
      default:
        break;
    }
  }

  addNotify(message, type = 'primary') {
    this.setState({
      notifications: this.state.notifications.add({ key: Date.now(), message, type }),
    });
  }

  handleCloseNotify(key) {
    this.setState({
      notifications: this.state.notifications.filter(n => n.key !== key),
    });
  }

  render() {
    const { metrics } = this.state;
    return (
      <div>
        <Header />
        <div className="container">
          <section className="hero">
            <div className="hero-body">
              <div className="container">
                <h1 className="title">
                  Trunks
                </h1>
                <h2 className="subtitle">
                  Trunks is a simple HTTP load testing tool with UI
                </h2>
                <FromPostAttack addNotify={this.addNotify} />
              </div>
            </div>
          </section>
          {metrics.is_active ? <Metrics metrics={metrics} /> : null}
        </div>
        <Notifications
          notifications={this.state.notifications.toArray()}
          onClose={this.handleCloseNotify}
        />
      </div>
    );
  }
}

export default App;