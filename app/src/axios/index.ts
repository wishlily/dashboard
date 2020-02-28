/**
 * Created by hao.cheng on 2017/4/16.
 */
import axios from 'axios';
import { get, post } from './tools';
import * as config from './config';
import moment from 'moment';

export const getBbcNews = () => get({ url: config.NEWS_BBC });

export const npmDependencies = () =>
    axios
        .get('./npm.json')
        .then(res => res.data)
        .catch(err => console.log(err));

export const weibo = () =>
    axios
        .get('./weibo.json')
        .then(res => res.data)
        .catch(err => console.log(err));

export const gitOauthLogin = () =>
    get({
        url: `${config.GIT_OAUTH}/authorize?client_id=792cdcd244e98dcd2dee&redirect_uri=http://localhost:3006/&scope=user&state=reactAdmin`,
    });
export const gitOauthToken = (code: string) =>
    post({
        url: `https://cors-anywhere.herokuapp.com/${config.GIT_OAUTH}/access_token`,
        data: {
            client_id: '792cdcd244e98dcd2dee',
            client_secret: '81c4ff9df390d482b7c8b214a55cf24bf1f53059',
            redirect_uri: 'http://localhost:3006/',
            state: 'reactAdmin',
            code,
        },
    });
// {headers: {Accept: 'application/json'}}
export const gitOauthInfo = (access_token: string) =>
    get({ url: `${config.GIT_USER}access_token=${access_token}` });

// easy-mock数据交互
// 管理员权限获取
export const admin = () => get({ url: config.MOCK_AUTH_ADMIN });
// 访问权限获取
export const guest = () => get({ url: config.MOCK_AUTH_VISITOR });

export const minVaildTS = () => moment(config.MIN_VAILD_TIME).unix();

// API
export interface IRecordParam {
    uuid: string;
    type: string;
    time: string;
    amount: number;
    account: Array<string>;
    unit?: number;
    nuv?: number;
    class?: Array<string>;
    member?: string;
    proj?: string;
    note?: string;
    deadline?: string;
}

export const getFinanceRecord = (t1: string, t2: string) => get({
    url: config.API_FINANCE_RECORD,
    config: {
        params: {
            start: t1,
            end: t2
        }
    }
});

export const setFinanceRecord = (type: string, data: IRecordParam) => post({
    url: config.API_FINANCE_RECORD,
    data: {
        type: type,
        data: data
    }
});

export interface IAccountParam {
    time: string;
    id: string;
    type: string;
    amount: number;
    unit?: number;
    nuv?: number;
    class?: string;
    deadline?: string;
    member?: string;
    account?: string;
    note?: string;
}

export const getFinanceAccount = (param?: string) => {
    if (param)
        return get({url: config.API_FINANCE_ACCOUNT + '?' + param});
    return get({url: config.API_FINANCE_ACCOUNT});
}

export const setFinanceAccount = (type: string, data: IAccountParam) => post({
    url: config.API_FINANCE_ACCOUNT,
    data: {
        type: type,
        data: data
    }
});