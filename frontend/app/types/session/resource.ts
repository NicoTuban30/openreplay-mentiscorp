import Record from 'Types/Record';
import { getResourceName } from 'App/utils';

const XHR = 'xhr' as const;
const FETCH = 'fetch' as const;
const JS = 'script' as const;
const CSS = 'css' as const;
const IMG = 'img' as const;
const MEDIA = 'media' as const;
const OTHER = 'other' as const;

function getResourceStatus(status: number, success: boolean) {
  if (status !== undefined) return String(status);
  if (typeof success === 'boolean' || typeof success === 'number') {
    return !!success
      ? '2xx-3xx'
      : '4xx-5xx';
  }
  return '2xx-3xx';
}

export const TYPES = {
  XHR,
  FETCH,
  JS,
  CSS,
  IMG,
  MEDIA,
  OTHER,
  "stylesheet": CSS,
}

const YELLOW_BOUND = 10;
const RED_BOUND = 80;

export function isRed(r: IResource) {
  return !r.success || r.score >= RED_BOUND;
}

interface IResource {
  type: keyof typeof TYPES,
  url: string,
  name: string,
  status: number,
  duration: number,
  index: number,
  time: number,
  ttfb: number,
  timewidth: number,
  success: boolean,
  score: number,
  method: string,
  request: string,
  response: string,
  headerSize: number,
  encodedBodySize: number,
  decodedBodySize: number,
  responseBodySize: number,
  timings: Record<string, any>
  datetime: number
  timestamp: number
}

export default class Resource {
  name = 'Resource'
  type: IResource["type"]
  status: string
  success: IResource["success"]
  time: IResource["time"]
  ttfb: IResource["ttfb"]
  url: IResource["url"]
  duration: IResource["duration"]
  index: IResource["index"]
  timewidth: IResource["timewidth"]
  score: IResource["score"]
  method: IResource["method"]
  request: IResource["request"]
  response: IResource["response"]
  headerSize: IResource["headerSize"]
  encodedBodySize: IResource["encodedBodySize"]
  decodedBodySize: IResource["decodedBodySize"]
  responseBodySize: IResource["responseBodySize"]
  timings: IResource["timings"]

  constructor({ status, success, time, datetime, timestamp, timings, ...resource }: IResource) {

    // adjusting for 201, 202 etc
    const reqSuccess = 300 > status || success
    Object.assign(this, {
      ...resource,
      name: getResourceName(resource.url),
      status: getResourceStatus(status, success),
      success: reqSuccess,
      time: typeof time === 'number' ? time : datetime || timestamp,
      ttfb: timings && timings.ttfb,
      timewidth: timings && timings.timewidth,
      timings,
      isRed: !reqSuccess || resource.score >= RED_BOUND,
      isYellow: resource.score < RED_BOUND && resource.score >= YELLOW_BOUND,
    })
  }
}

