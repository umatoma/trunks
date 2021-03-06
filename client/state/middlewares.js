import debug from 'debug';

const debugState = debug('trunks:state');

/**
 * logging middleware
 */
const loggerMiddleware = (state, actionName) => {
  debugState('%s %O', actionName, state);
  return state;
};

/**
 * get middlewares
 * @return {Array.<{Function}>}
 */
export function getMiddlewares() { // eslint-disable-line
  const middlewares = [];
  if (process.env.NODE_ENV !== 'production') {
    middlewares.push(loggerMiddleware);
  }
  return middlewares;
}
