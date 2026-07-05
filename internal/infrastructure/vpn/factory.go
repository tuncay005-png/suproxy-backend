package vpn

import (
	"errors"
	"fmt"
)

var (
	ErrUnsupportedKernel = errors.New("unsupported kernel type")
	ErrKernelNotFound    = errors.New("kernel not found")
)

// kernelFactory implements KernelFactory
type kernelFactory struct {
	kernels map[string]Kernel
}

// NewKernelFactory creates a new kernel factory
func NewKernelFactory() KernelFactory {
	return &kernelFactory{
		kernels: make(map[string]Kernel),
	}
}

// Register registers a kernel implementation
func (f *kernelFactory) Register(kernel Kernel) {
	f.kernels[kernel.Name()] = kernel
}

func (f *kernelFactory) Create(kernelType string) (Kernel, error) {
	kernel, ok := f.kernels[kernelType]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedKernel, kernelType)
	}
	return kernel, nil
}

func (f *kernelFactory) SupportedKernels() []string {
	kernels := make([]string, 0, len(f.kernels))
	for name := range f.kernels {
		kernels = append(kernels, name)
	}
	return kernels
}

// DefaultFactory creates a factory with default kernels
// This will be used in bootstrap/dependency injection
func DefaultFactory(xrayKernel Kernel) KernelFactory {
	factory := NewKernelFactory().(*kernelFactory)
	factory.Register(xrayKernel)
	// Future: factory.Register(singboxKernel)
	// Future: factory.Register(hysteriaKernel)
	return factory
}
