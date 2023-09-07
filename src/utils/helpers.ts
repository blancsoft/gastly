import zlib from "zlib";

export const brotli = {
  getParams(mode?: number, quality?: number, sizeHint?: number) {
    const params =  {
      [zlib.constants.BROTLI_PARAM_MODE]: mode ?? zlib.constants.BROTLI_DEFAULT_MODE,
      [zlib.constants.BROTLI_PARAM_QUALITY]: quality ?? zlib.constants.BROTLI_MAX_QUALITY
    }
    if (sizeHint) {
      params[zlib.constants.BROTLI_PARAM_SIZE_HINT] = sizeHint
    }
    return params
  },

  compress(content: Buffer, mode?: number, quality?: number) {
    const params = this.getParams(mode, quality, content.byteLength)
    return new Promise<Buffer>((resolve, reject) => {
      zlib.brotliCompress(content, { params }, (err, result) =>
        err ? reject(err) : resolve(result)
      )
    })
  },

  decompress(content: Buffer, mode?: number, quality?: number) {
    const params = this.getParams(mode, quality, content.byteLength)
    return new Promise<Buffer>((resolve, reject) => {
      zlib.brotliDecompress(content, { params }, (err, result) =>
        err ? reject(err) : resolve(result)
      )
    })
  }
}

export default brotli;
