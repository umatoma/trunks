import React from 'react';
import { Link } from 'react-router-dom';
import FromPostAttack from './FormPostAttack';
import WorkerProgress from './WorkerProgress';
import Metrics from './Metrics';
import ModalImportOption from './ModalImportOption';
import { convertToObject } from '../lib/option';

const isActive = worker => worker.status === 'active';
const isDone = worker => worker.status === 'done';
const isNotReady = worker => worker.status !== 'ready';

const Attack = ({
  worker,
  metrics,
  formAttack,
  importOption,
  addNotify,
  updateFormAttack,
  setFormAttack,
  updateFormImport,
  openImportModal,
  closeImportModal,
}) => (
  <div>
    <div className="content">
      <div className="message">
        <div className="message-body">
          <span className="icon">
            <i className="fa fa-lightbulb-o" />
          </span>
          <span>
            Refer to <a href="https://github.com/tsenart/vegeta" target="_blank" rel="noopener noreferrer">vegeta</a> for explanation of each option.
          </span>
        </div>
      </div>
      <div className="field is-grouped" style={{ marginBottom: '1rem' }}>
        <p className="control">
          <button
            className="button is-small"
            onClick={openImportModal}
          >
            <span className="icon is-small">
              <i className="fa fa-upload" />
            </span>
            <span>Import options</span>
          </button>
        </p>
      </div>
      <FromPostAttack
        form={formAttack}
        addNotify={addNotify}
        onUpdate={updateFormAttack}
        isAttacking={worker.status === 'active'}
      />
    </div>
    <div className="content">
      {isNotReady(worker, metrics) ? <Metrics metrics={metrics} /> : null}
      {isActive(worker) ? <WorkerProgress worker={worker} metrics={metrics} /> : null}
      {isDone(worker) ? (
        <Link className="button is-primary is-outlined is-fullwidth" to={`/results/${worker.filename}`}>
          Check details
        </Link>
      ) : null}
    </div>
    <ModalImportOption
      isActive={importOption.isModalActive}
      form={{ text: importOption.text, error: importOption.error }}
      onChange={updateFormImport}
      onClose={closeImportModal}
      onSubmit={({ text }) => {
        try {
          setFormAttack(convertToObject(text));
          closeImportModal();
          addNotify('Importing option was successful.');
        } catch (e) {
          updateFormImport({ error: e });
        }
      }}
    />
  </div>
);

Attack.propTypes = {
  worker: React.PropTypes.object.isRequired,
  metrics: React.PropTypes.object.isRequired,
  formAttack: React.PropTypes.object.isRequired,
  importOption: React.PropTypes.object.isRequired,
  addNotify: React.PropTypes.func.isRequired,
  updateFormAttack: React.PropTypes.func.isRequired,
  setFormAttack: React.PropTypes.func.isRequired,
  updateFormImport: React.PropTypes.func.isRequired,
  openImportModal: React.PropTypes.func.isRequired,
  closeImportModal: React.PropTypes.func.isRequired,
};

export default Attack;
