import React from 'react';
import { Link } from 'react-router-dom';
import FromPostAttack from './FormPostAttack';
import WorkerProgress from './WorkerProgress';
import Metrics from './Metrics';

const isActive = worker => worker.status === 'active';
const isDone = worker => worker.status === 'done';
const isActiveOrDone = worker => isActive(worker) || isDone(worker);

const Attack = ({ worker, metrics, addNotify }) => (
  <div>
    <section className="section">
      <FromPostAttack addNotify={addNotify} isAttacking={worker.status === 'active'} />
    </section>
    <section className="section">
      {isActiveOrDone(worker, metrics) ? <Metrics metrics={metrics} /> : null}
      {isActive(worker) ? <WorkerProgress worker={worker} metrics={metrics} /> : null}
      {isDone(worker) ? (
        <Link className="button is-primary is-outlined is-fullwidth" to={`/results/${worker.filename}`}>
          Check details
        </Link>
      ) : null}
    </section>
  </div>
);

Attack.propTypes = {
  worker: React.PropTypes.object.isRequired,
  metrics: React.PropTypes.object.isRequired,
  addNotify: React.PropTypes.func.isRequired,
};

export default Attack;
