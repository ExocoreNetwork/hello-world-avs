package chainio

import (
	sdkmath "cosmossdk.io/math"
	avs "github.com/ExocoreNetwork/exocore-avs/contracts/bindings/avs"
	"github.com/ExocoreNetwork/exocore-avs/core/chainio/eth"
	"github.com/ExocoreNetwork/exocore-sdk/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type ExoReader interface {
	GetOptInOperators(
		opts *bind.CallOpts,
		avsAddress string,
	) ([]string, error)

	GetRegisteredPubkey(
		opts *bind.CallOpts,
		operator string,
	) ([]byte, error)
	GtAVSUSDValue(
		opts *bind.CallOpts,
		avsAddress string,
	) (sdkmath.LegacyDec, error)

	GetOperatorOptedUSDValue(
		opts *bind.CallOpts,
		avsAddress string,
		operatorAddr string,
	) (sdkmath.LegacyDec, error)
	GetAVSEpochIdentifier(
		opts *bind.CallOpts,
		avsAddress string,
	) (string, error)
	GetTaskInfo(
		opts *bind.CallOpts,
		avsAddress string,
		taskID uint64,
	) ([]uint64, error)
	IsOperator(
		opts *bind.CallOpts,
		operator string,
	) (bool, error)

	GetCurrentEpoch(
		opts *bind.CallOpts,
		epochIdentifier string,
	) (int64, error)
}

type ExoChainReader struct {
	logger     logging.Logger
	avsManager avs.ContracthelloWorld
	ethClient  eth.EthClient
}

// forces EthReader to implement the chainio.Reader interface
var _ ExoReader = (*ExoChainReader)(nil)

func NewExoChainReader(
	avsManager avs.ContracthelloWorld,
	logger logging.Logger,
	ethClient eth.EthClient,
) *ExoChainReader {
	return &ExoChainReader{
		avsManager: avsManager,
		logger:     logger,
		ethClient:  ethClient,
	}
}

func BuildExoChainReader(
	avsAddr gethcommon.Address,
	ethClient eth.EthClient,
	logger logging.Logger,
) (*ExoChainReader, error) {
	exoContractBindings, err := NewExocoreContractBindings(
		avsAddr,
		ethClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	return NewExoChainReader(
		*exoContractBindings.AVSManager,
		logger,
		ethClient,
	), nil
}

func (r *ExoChainReader) GetOptInOperators(
	opts *bind.CallOpts,
	avsAddress string,
) ([]string, error) {
	operators, err := r.avsManager.GetOptInOperators(
		opts,
		gethcommon.HexToAddress(avsAddress))
	if err != nil {
		r.logger.Error("Failed to GetOptInOperators ", "err", err)
		return nil, err
	}
	return operators, nil
}

func (r *ExoChainReader) GetRegisteredPubkey(opts *bind.CallOpts, operator string) ([]byte, error) {
	pukKey, err := r.avsManager.GetRegisteredPubkey(
		opts,
		operator)
	if err != nil {
		r.logger.Error("Failed to GetRegisteredPubkey ", "err", err)
		return nil, err
	}
	return pukKey, nil
}

func (r *ExoChainReader) GtAVSUSDValue(opts *bind.CallOpts, avsAddress string) (sdkmath.LegacyDec, error) {
	amount, err := r.avsManager.GetAVSUSDValue(
		opts,
		gethcommon.HexToAddress(avsAddress))
	if err != nil {
		r.logger.Error("Failed to GtAVSUSDValue ", "err", err)
		return sdkmath.LegacyDec{}, err
	}
	return sdkmath.LegacyNewDecFromBigInt(amount), nil
}

func (r *ExoChainReader) GetOperatorOptedUSDValue(opts *bind.CallOpts, avsAddress string, operatorAddr string) (sdkmath.LegacyDec, error) {
	amount, err := r.avsManager.GetOperatorOptedUSDValue(
		opts,
		gethcommon.HexToAddress(avsAddress), operatorAddr)
	if err != nil {
		r.logger.Error("Failed to GetOperatorOptedUSDValue ", "err", err)
		return sdkmath.LegacyDec{}, err
	}
	return sdkmath.LegacyNewDecFromBigInt(amount), nil
}

func (r *ExoChainReader) GetAVSEpochIdentifier(opts *bind.CallOpts, avsAddress string) (string, error) {
	epochIdentifier, err := r.avsManager.GetAVSEpochIdentifier(
		opts,
		gethcommon.HexToAddress(avsAddress))
	if err != nil {
		r.logger.Error("Failed to GetAVSEpochIdentifier ", "err", err)
		return "", err
	}
	return epochIdentifier, nil
}
func (r *ExoChainReader) GetTaskInfo(opts *bind.CallOpts, avsAddress string, taskID uint64) ([]uint64, error) {
	info, err := r.avsManager.GetTaskInfo(
		opts,
		gethcommon.HexToAddress(avsAddress), taskID)
	if err != nil {
		r.logger.Error("Failed to GetTaskInfo ", "err", err)
		return nil, err
	}
	return info, nil
}

func (r *ExoChainReader) IsOperator(opts *bind.CallOpts, operator string) (bool, error) {
	flag, err := r.avsManager.IsOperator(
		opts,
		gethcommon.HexToAddress(operator))
	if err != nil {
		r.logger.Error("Failed to exec IsOperator ", "err", err)
		return false, err
	}
	return flag, nil
}
func (r *ExoChainReader) GetCurrentEpoch(opts *bind.CallOpts, epochIdentifier string) (int64, error) {
	currentEpoch, err := r.avsManager.GetCurrentEpoch(
		opts,
		epochIdentifier)
	if err != nil {
		r.logger.Error("Failed to exec IsOperator ", "err", err)
		return 0, err
	}
	return currentEpoch, nil
}
