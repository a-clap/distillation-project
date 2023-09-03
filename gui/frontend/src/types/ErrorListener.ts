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

import { NotifyError } from '../../wailsjs/go/backend/Events'
import { Listener } from './Listener';
import { AppErrorCodes } from '../stores/error_codes';

declare type ErrorCallback = (id: number) => void;

class errorListener {
  private static _instance: errorListener;
  error: Listener;

  public static get Instance() {
    return this._instance || (this._instance = new this());
  }

  private constructor() {
    this.error = new Listener()

    NotifyError().then((ev: string) => {
      return runtime.EventsOn(ev, (...args: any) => {
        this.handle(...args);
      });
    })

  }

  private handle(...args: any) {
    try {
      let t = Number(args[0])
      this.error.notify(t)
    } catch (e) {
      console.log(e)
    }
  }

  sendError(nb: AppErrorCodes) {
    this.handle(nb)
  }


  subscribe(cb: ErrorCallback) {
    this.error.subscribe(cb)
  }

  unsubscribe(cb: ErrorCallback) {
    this.error.unsubscribe(cb)
  }

}

export var ErrorListener = errorListener.Instance
