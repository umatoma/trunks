import { Record, List, OrderedSet, Map } from 'immutable';

export const ModelHeader = Record({
  isHamburgerActive: false,
});

export const ModelSideMenu = Record({
  isModalActive: false,
});

export const ModelImportOption = Record({
  isModalActive: false,
  error: null,
  text: '',
});

export const ModelWorker = Record({
  status: 'ready',
  error: null,
  duration: 0,
  rate: 0,
  filename: '',
});

export const ModelMetrics = Record({
  bytes_in: { total: 0, mean: 0 },
  bytes_out: { total: 0, mean: 0 },
  duration: 0,
  earliest: '',
  end: '',
  errors: null,
  latencies: { mean: 0, max: 0, '50th': 0, '95th': 0, '99th': 0 },
  latest: {},
  rate: 0,
  requests: 0,
  status_codes: {},
  success: 0,
  wait: 0,
});

export const ModelReport = Record({
  isFetching: true,
  showResultList: false,
  error: null,
  metrics: new ModelMetrics(),
  histgram: List(),
  results: List(),
});

export const ModelFormAttack = Record({
  Body: '',
  Duration: '10s',
  Rate: 1,
  Targets: '',
});

export const getInitialState = state => Object.assign({
  notifications: OrderedSet(),
  header: new ModelHeader(),
  sideMenu: new ModelSideMenu(),
  importOption: new ModelImportOption(),
  worker: new ModelWorker(),
  metrics: new ModelMetrics(),
  resultFiles: List(),
  reports: Map(),
  formAttack: new ModelFormAttack(),
}, state);