export interface IFMenuBase {
    key: string;
    title: string;
    icon?: string;
    component?: string;
    query?: string;
    auth?: string;
    route?: string;
    login?: string;
}

export interface IFMenu extends IFMenuBase {
    subs?: IFMenuBase[];
}

const menus: {
    menus: IFMenu[];
    others: IFMenu[] | [];
    [index: string]: any;
} = {
    menus: [
        // 菜单相关路由
        { key: '/app/dashboard/index', title: '首页', icon: 'home', component: 'Dashboard' },
        {
            key: '/app/finance',
            title: '财务',
            icon: 'wallet',
            subs: [
                // TODO: not ok
                { key: '/app/finance/record', title: '流水', component: 'ViewRecords' },
                { key: '/app/finance/account', title: '账户', component: '' },
            ],
        },
    ],
    others: [], // 非菜单相关路由
};

export default menus;
