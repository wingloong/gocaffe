package layer

import (
	"errors"
	"log"
	"math"

	pb "github.com/cvley/gocaffe"
	"github.com/cvley/gocaffe/blob"
	"github.com/cvley/gocaffe/math"
)

type ConvLayer struct {
	ConvParam *pb.ConvolutionParameter
	IsShared  bool
	Phase     *pb.Phase
	blobs     []*blob.Blob

	forceNdim2col   bool
	channelAxis     int
	numSpatialAxis  int
	kernelShapeData *blob.Blob
	stride          *blob.Blob
	pad             *blob.Blob
}

func NewConvolutionLayer(param *pb.LayerParameter) *ConvLayer {
	phase := param.GetPhase()
	blobprotos := param.GetBlobs()
	blobs := make([]*blob.Blob, len(blobprotos))
	for i, p := range blobprotos {
		blobs[i].FromProto(p, true)
	}
	return &ConvLayer{
		ConvParam: param.GetConvolutionParam(),
		IsShared:  false,
		Phase:     phase,
		blobs:     blobs,
	}
}

func (conv *ConvLayer) SetUp(bottom, top []*blob.Blob) error {
	conv.forceNdim2col = conv.ConvParam.GetForceNdIm2Col()
	conv.channelAxis = bottom[0].CanonicalAxisIndex(int(conv.ConvParam.GetAxis()))

	firstSpaticalAxis := conv.channelAxis + 1
	numAxes := bottom[0].AxesNum()
	conv.numSpatialAxis = numAxes - firstSpaticalAxis
	if conv.numSpaticalAxis < 0 {
		return errors.New("conv layer num spatial axis less than 0")
	}

	bottomDimBlobShape := make([]int32, conv.numSpatialAxis+1)
	spatialDimBlobShape := make([]int32, 1)
	if conv.numSpatialAxis > 1 {
		spatialDimBlobShape = make([]int32, conv.numSpatialAxis)
	}

	// setup filter kernel dimensions (kernel_shape)
	conv.kernel_shape_data.Reshape(spatialDimBlobShape)
	if conv.ConvParam.GetKernelH() > 0 || conv.ConvParam.GetKernelW() > 0 {
		if conv.numSpatialAxis != 2 {
			return errors.New("kernel_h & kernel_w can only be used for 2D convolution.")
		}
		if len(conv.ConvParam.GetKernelSize()) != 0 {
			return errors.New("Either kernel_size or kernel_h/w should be specified; not both.")
		}
		conv.kernelShapeData.Data[0] = conv.ConvParam.GetKernelH()
		conv.kernelShapeData.Data[1] = conv.ConvParam.GetKernelW()
	} else {
		numKernelDims := len(conv.ConvParam.GetKernelSize())
		if numKernelDims == 1 || numKernelDims == conv.numSpatialAxis {
			return errors.New("kernel_size must be specified once, or once per spatial dimension.")
		}
		for i := 0; i < conv.numSpatialAxis; i++ {
			if numKernelDims == 1 {
				conv.kernelShapeData.Data[i] = conv.ConvParam.GetKernelSize()[0]
			} else {
				conv.kernelShapeData.Data[i] = conv.ConvParam.GetKernelSize()[i]
			}
		}
	}

	for i = 0; i < conv.numSpatialAxis; i++ {
		if conv.kernelShapeData[i] <= 0 {
			return errors.New("Filter dimensions must be nonzeros.")
		}
	}

	// setup stride dimensions
	conv.stride.Reshape(spatialDimBlobShape)
	if conv.ConvParam.GetStrideH() > 0 || conv.ConvParam.GetStrideW() > 0 {
		if conv.numSpatialAxis != 2 {
			return errors.New("stride_h & stride_w can only be used for 2D convolution.")
		}
		if len(conv.ConvParam.GetStride()) != 0 {
			return errors.New("Either stride or stride_h/w should be specified; not both.")
		}
		conv.stride.Data[0] = conv.ConvParam.GetStrideH()
		conv.stride.Data[1] = conv.ConvParam.GetStrideW()
	} else {
		numStrideDims := len(conv.ConvParam.GetStride())
		if numStrideDims == 0 || numStrideDims == 1 || numStrideDims == conv.numSpatialAxis {
			return errors.New("stride must be specified once, or once per spatical dimension.")
		}
		kDefaultStride := 1
		for i := 0; i < conv.numSpatialAxis; i++ {
			if numStrideDims == 0 {
				conv.stride.Data[i] = kDefaultStride
			} else if numStrideDims == 1 {
				conv.stride.Data[i] = conv.ConvParam.GetStride()[0]
			} else {
				conv.stride.Data[i] = conv.ConvParam.GetStride()[i]
			}
			if conv.stride.Data[i] <= 0 {
				return errors.New("stride dimensions must be nonzero.")
			}
		}
	}

	// setup pad dimensions
	conv.pad.Reshape(spatialDimBlobShape)
	if conv.ConvParam.GetPadH() > 0 || conv.ConvParam.GetPadW() > 0 {
		if conv.numSpatialAxis != 2 {
			return errors.New("pad_h & pad_w can only be used for 2D convolution.")
		}
		if len(conv.ConvParam.GetPad() != 0) {
			return errors.New("Either pad or pad_h/w should be specified; not both.")
		}
		conv.pad.Data[0] = conv.ConvParam.GetPadH()
		conv.pad.Data[1] = conv.ConvParam.GetPadW()
	} else {
		numPadDims := len(conv.ConvParam.GetPad())
		if numPadDims != 0 && numPadDims != 1 && numPadDims != conv.numSpatialAxis {
			return errors.New("pad must be specified once, or once per spatial dimension")
		}
		kDefaultPad := 0
		for i := 0; i < conv.numSpatialAxis; i++ {
			switch numPadDims {
			case 0:
				conv.pad.Data[i] = kDefaultPad
			case 1:
				conv.pad.Data[i] = conv.ConvParam.GetPad()[0]
			default:
				conv.pad.Data[i] = conv.ConvParam.GetPad()[i]
			}
		}
	}

	// setup dilation dimensions
}

func (conv *ConvLayer) Reshape(bottom, top []*blob.Blob) {
}

func (conv *ConvLayer) Forward(bottom, top []*blob.Blob) {
}

func (conv *ConvLayer) Backward(bottom, top []*blob.Blob, propagateDown []bool) {
	// not implement yet, only forward is enough
}

func (conv *ConvLayer) Type() string {
	return "ConvolutionLayer"
}

func im2colNd() {

}

func im2col() {
}
