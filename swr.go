/*
2015  Sleepy Programmer <hunan@emsym.com>
*/
package gmf

/*

#cgo pkg-config: libswresample

#include "libswresample/swresample.h"
#include <libavcodec/avcodec.h>
#include <libavutil/frame.h>

int gmf_swr_convert(SwrContext *ctx, int srcSampleRate, AVFrame *srcFrame, int dstSampleRate, AVFrame *dstFrame){
    int nbDstSamples = dstFrame->nb_samples;

    if (srcSampleRate && dstSampleRate) {
        nbDstSamples = av_rescale_rnd(swr_get_delay(ctx, srcSampleRate) + srcFrame->nb_samples, dstSampleRate, srcSampleRate, AV_ROUND_UP);
    }

	int nbSamplesConverted = swr_convert(ctx, dstFrame->data, nbDstSamples, (const uint8_t **)srcFrame->data, srcFrame->nb_samples);
    if (nbSamplesConverted < 0) {
        return nbSamplesConverted;
    }

    return av_samples_get_buffer_size(0, dstFrame->channel_layout, nbSamplesConverted, dstFrame->format, 1);
}

*/
import "C"

type SwrCtx struct {
	srcSampleRate C.int
	dstSampleRate C.int
	swrCtx        *C.struct_SwrContext
	CgoMemoryManage
}

func NewSwrCtx(options []*Option) *SwrCtx {
	this := &SwrCtx{swrCtx: C.swr_alloc()}

	for _, option := range options {
		if option.Key == "in_sample_rate" {
			if srcSampleRate, ok := option.Val.(int); ok {
				this.srcSampleRate = C.int(srcSampleRate)
			}
		}

		if option.Key == "out_sample_rate" {
			if dstSampleRate, ok := option.Val.(int); ok {
				this.dstSampleRate = C.int(dstSampleRate)
			}
		}

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
	samplesBufferSize := C.gmf_swr_convert(this.swrCtx, this.srcSampleRate, src.avFrame, this.dstSampleRate, dst.avFrame)

	return int(samplesBufferSize)
}
