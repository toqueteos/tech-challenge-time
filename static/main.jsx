const defaultFetchOptions = {
  mode: "cors",
  cache: "no-cache",
  headers: {
    "Content-Type": "application/json"
  },
  referrerPolicy: "no-referrer"
};

const getActiveSessions = async () => {
  const url = "/timer/active";
  const options = { ...defaultFetchOptions, method: "GET" };
  const response = await fetch(url, options);
  return await response.json();
};

const getSessionsOf = async period => {
  const url = `/timer/analytics?period=${period}`;
  const options = { ...defaultFetchOptions, method: "GET" };
  const response = await fetch(url, options);
  return await response.json();
};

const getSessionsOfDay = async () => await getSessionsOf("day");
const getSessionsOfWeek = async () => await getSessionsOf("week");
const getSessionsOfMonth = async () => await getSessionsOf("month");

const startSession = async name => {
  const url = "/timer/start";
  const options = {
    ...defaultFetchOptions,
    method: "POST",
    body: JSON.stringify({ name })
  };
  const response = await fetch(url, options);
  return await response.json();
};

const stopSession = async sessionId => {
  const url = "/timer/stop";
  const options = {
    ...defaultFetchOptions,
    method: "POST",
    body: JSON.stringify({ id: sessionId })
  };
  const response = await fetch(url, options);
  return await response.json();
};

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = { page: "active" };

    this.updatePage = this.updatePage.bind(this);
  }

  updatePage(page) {
    this.setState({ page });
  }

  render() {
    const pages = ["active", "day", "week", "month"];
    const pageNotFound = !pages.includes(this.state.page);

    return (
      <div className="container">
        <Actions active={this.state.page} eventHandler={this.updatePage} />
        {this.state.page === "active" ? <ActiveSessions /> : null}
        {this.state.page === "day" ? (
          <SessionsOf title="Today's Sessions" refresh={getSessionsOfDay} />
        ) : null}
        {this.state.page === "week" ? (
          <SessionsOf
            title="This Week's Sessions"
            refresh={getSessionsOfWeek}
          />
        ) : null}
        {this.state.page === "month" ? (
          <SessionsOf
            title="This Month's Sessions"
            refresh={getSessionsOfMonth}
          />
        ) : null}
        {pageNotFound ? <PageNotFound /> : null}
      </div>
    );
  }
}

function PageNotFound(props) {
  return (
    <div className="row mt-3">
      <div className="col">
        <h1>404 Page not found</h1>
      </div>
    </div>
  );
}

function Actions(props) {
  const onClickFn = props.eventHandler;
  const buttonTexts = ["Active", "Day", "Week", "Month"];
  const buttons = buttonTexts.map(text => {
    const key = text.toLowerCase();
    const classNames =
      props.active === key ? "btn btn-primary" : "btn btn-secondary";
    return (
      <button
        key={key}
        className={classNames}
        onClick={onClickFn.bind(this, key)}
      >
        {text}
      </button>
    );
  });

  return (
    <div className="row">
      <div className="col">
        <div
          className="btn-group"
          role="group"
          aria-label="My interpretation of a silly yet functional navbar/router"
        >
          {buttons}
        </div>
      </div>
    </div>
  );
}

class SessionsOf extends React.Component {
  constructor(props) {
    super(props);
    this.state = { sessions: [] };

    this.refresh = this.refresh.bind(this);
  }

  async componentDidMount() {
    await this.refresh();
  }

  async refresh() {
    this.setState({
      sessions: await this.props.refresh()
    });
  }

  render() {
    return (
      <div className="row mt-3">
        <div className="col">
          <SessionList
            title={this.props.title}
            sessions={this.state.sessions}
            refresh={this.refresh}
          />
        </div>
      </div>
    );
  }
}

class ActiveSessions extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      sessions: []
    };

    this.refresh = this.refresh.bind(this);
  }

  async componentDidMount() {
    await this.refresh();
  }

  async refresh() {
    this.setState({
      sessions: await getActiveSessions()
    });
  }

  render() {
    return (
      <div className="row mt-3 mb-3">
        <div className="col-sm-4">
          <NewSession refresh={this.refresh} />
        </div>
        <div className="col-sm-8">
          <SessionList
            title="Active Sessions"
            sessions={this.state.sessions}
            refresh={this.refresh}
          />
        </div>
      </div>
    );
  }
}

class NewSession extends React.Component {
  constructor(props) {
    super(props);
    this.state = { name: "" };

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleClick = this.handleClick.bind(this);
    this.handleSessionStart = this.handleSessionStart.bind(this);
  }

  handleChange(event) {
    this.setState({ name: event.target.value });
  }

  async handleSubmit(event) {
    event.preventDefault();
    await this.handleSessionStart();
  }

  async handleClick(event) {
    event.preventDefault();
    await this.handleSessionStart();
  }

  async handleSessionStart() {
    const name = this.state.name;
    this.setState({ name: "" });
    await startSession(name);
    await this.props.refresh();
  }

  render() {
    return (
      <div className="card">
        <h5 className="card-header">New Session</h5>
        <div className="card-body">
          <div className="card-text">
            <form onSubmit={this.handleSubmit}>
              <div className="form-group">
                <label htmlFor="name">Session name</label>
                <input
                  id="input-name"
                  type="text"
                  className="form-control"
                  placeholder="Ex: Review Emil's PR, Lunch with Jonas..."
                  value={this.state.name}
                  onChange={this.handleChange}
                />
              </div>
              <div className="form-group">
                <button
                  id="button-track"
                  type="button"
                  className="btn btn-primary"
                  onClick={this.handleClick}
                >
                  <i className="far fa-clock"></i>
                  Track!
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    );
  }
}

function SessionList(props) {
  return (
    <div className="card">
      <h5 className="card-header">{props.title}</h5>
      <div className="card-body">
        <SessionTable sessions={props.sessions} refresh={props.refresh} />
      </div>
    </div>
  );
}

function SessionTable(props) {
  let rows;

  if (!props.sessions || props.sessions.length === 0) {
    rows = (
      <tr>
        <td colSpan="5" className="text-center">
          No sessions found. Why don't you start one?
        </td>
      </tr>
    );
  } else {
    rows = props.sessions.map(session => (
      <SessionRow key={session.id} session={session} refresh={props.refresh} />
    ));
  }

  return (
    <table className="table">
      <thead>
        <tr>
          <th scope="col">Name</th>
          <th scope="col">Started</th>
          <th scope="col">Ended</th>
          <th scope="col">Elapsed</th>
          <th scope="col">Actions</th>
        </tr>
      </thead>
      <tbody>{rows}</tbody>
    </table>
  );
}

class SessionRow extends React.Component {
  constructor(props) {
    super(props);

    this.onStop = this.onStop.bind(this);
  }

  // Clicking the "Stop button" does three things:
  // - Performs an XHR to backend (to stop this session)
  // - Changes state to hide this session from DOM
  // - Forces a data reload (SessionTable gets refreshed)
  async onStop(sessionId) {
    await stopSession(sessionId);
    await this.props.refresh();
  }

  render() {
    const fmt = "YYYY/MM/DD HH:mm:ss";
    const session = this.props.session;

    const start = moment.utc(session.date_start);
    const end = moment.utc(session.date_end);
    const diff = moment.duration(end.diff(start));
    const diffNow = moment.duration(moment.utc().diff(start));

    return (
      <tr>
        <td className="text-truncate">{session.name}</td>
        <td>
          <abbr title={start.format(fmt)}>{start.fromNow()}</abbr>
        </td>
        <td>
          {session.date_end ? (
            <abbr title={end.format(fmt)}>{end.fromNow()}</abbr>
          ) : (
            "in progress"
          )}
        </td>
        <td>
          {session.date_end ? (
            <abbr title={`${diff.asHours()} hours`}>{diff.humanize()}</abbr>
          ) : (
            <abbr title={`${diffNow.asHours()} hours`}>
              {diffNow.humanize()}
            </abbr>
          )}
        </td>
        <td>
          {!session.date_end ? (
            <button onClick={this.onStop.bind(this, session.id)}>Stop</button>
          ) : null}
        </td>
      </tr>
    );
  }
}

ReactDOM.render(<App />, document.getElementById("root"));
