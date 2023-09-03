// MIT License
//
// Copyright (c) 2023 a-clap
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import {ElLoading} from "element-plus";
import {i18n} from "../main";
import {AppErrorCodes} from "../stores/error_codes";
import {ErrorListener} from "./ErrorListener";

let defaultText = i18n.global.t('loader.loading')

class loader {
    private loading: typeof ElLoading.service
    private timeoutID: number
    private errCode: AppErrorCodes

    constructor() {
        this.loading = ElLoading.service
        this.timeoutID = 0
        this.errCode = AppErrorCodes.Success
    }

    show(errCode: AppErrorCodes, timeout: number, text: string = defaultText) {
        this.loading = ElLoading.service({
            lock: true,
            text: text,
            background: 'rgba(0, 0, 0, 0.7)',
        })

        this.errCode = errCode

        this.timeoutID = setTimeout(() => {
            this.fail()
        }, timeout)
    }
    failNow(errCode: AppErrorCodes) {
        this.errCode = errCode
        this.fail()
    }
    fail() {
        ErrorListener.sendError(this.errCode)
        this.close()
    }


    close() {
        if (this.loading != null) {
            clearTimeout(this.timeoutID)
            this.loading.close()
        }
    }
}

export var Loader: loader = new loader()
