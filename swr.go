/*
2015  Sleepy Programmer <hunan@emsym.com>
*/
package gmf

/*

#cgo pkg-config: libswresample

#include "libswresample/swresample.h"
#include <libavcodec/avcodec.h>
#include <libavutil/frame.h>

int gmf_swr_convert(SwrContext *ctx, AVFrame *srcFrame, AVFrame *dstFrame){
	int nbSamplesConverted = swr_convert(ctx, dstFrame->data, dstFrame->nb_samples, (const uint8_t **)srcFrame->data, srcFrame->nb_samples);
    if (nbSamplesConverted < 0) {
        return nbSamplesConverted;
    }

    return av_samples_get_buffer_size(0, dstFrame->channel_layout, nbSamplesConverted, dstFrame->format, 0);
}

*/
import "C"

type SwrCtx struct {
	swrCtx *C.struct_SwrContext
	CgoMemoryManage
}

func NewSwrCtx(options []*Option) *SwrCtx {
	this := &SwrCtx{swrCtx: C.swr_alloc()}

	for _, option := range options {
		option.Set(this.swrCtx)
	}

	if int(C.swr_init(this.swrCtx)) < 0 {
		return nil
	}

	return this
}

func (this *SwrCtx) Free() {
	C.swr_free(&this.swrCtx)
}

func (this *SwrCtx) Convert(src *Frame, dst *Frame) int {
	samplesBufferSize := C.gmf_swr_convert(this.swrCtx, src.avFrame, dst.avFrame)

	return int(samplesBufferSize)
}
