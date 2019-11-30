import React, {Component} from 'react';
import {
    Table,
    Button,
    Row,
    Col,
    Icon
} from 'antd';
import {IAccountParam} from '../../axios'
import AccountForm from './AccountForm';

interface IAccountTable {
    key: string;
    type: string;
    amount: number;
    data: Array<IAccountParam>;
}

type AccountTableProps = {
    onChange?: (v: IAccountParam) => void;
    onDelete?: (v: IAccountParam) => void;
    onCreate?: (v: IAccountParam) => void;
    data: Array<IAccountParam>;
}
interface AccountTableState {
    amount: number;
    raw: Array<IAccountParam>;
    data: Array<IAccountTable>;
}

class AccountTable extends Component<AccountTableProps, AccountTableState> {
    constructor(props: any) {
        super(props);
        this.state = {
            raw: [],
            amount: 0,
            data: [],
        };
        this.columns = [{
            title: '类型',
            dataIndex: 'type',
        }, {
            title: '金额',
            dataIndex: 'amount',
        }];
    }
    static getDerivedStateFromProps(prevProps: AccountTableProps, prevState: AccountTableState) {
        const {data} = prevProps
        if (data && data !== prevState.raw) {
            const groups: Array<string> = []
            data.forEach(item => {
                if (item.class && (groups.indexOf(item.class) < 0) && (item.type.length > 1)) {
                    groups.push(item.class)
                }
            })
            const val: Array<IAccountTable> = []
            var amount = 0
            groups.forEach(type => {
                const item: IAccountTable = {key: type, type: type, amount: 0, data: []}
                data.forEach(v => {
                    if (v.class === type) {
                        var input = v.amount
                        if (v.unit && v.nuv && (v.nuv * v.unit) !== 0) {
                            input = v.nuv * v.unit
                        }
                        amount = +Number(amount + input).toFixed(2)
                        item.amount = +Number(item.amount + input).toFixed(2)
                        var t:any = v
                        t.key = v.id
                        item.data.push(t)
                    }
                })
                val.push(item)
            })
            return {
                raw: data,
                data: val,
                amount: amount
            };
        }
        return null
    }
    columns: any;
    expandedRowRender(record: IAccountTable) {
        const columns = [
            {
                title: '账户',
                dataIndex: 'id'
            }, {
                title: '币种',
                dataIndex: 'type'
            }, {
                title: '份额',
                dataIndex: 'unit'
            }, {
                title: '净值',
                dataIndex: 'nuv'
            }, {
                title: '投入',
                dataIndex: 'amount'
            }, {
                title: '到期',
                dataIndex: 'deadline'
            }, {
                title: '操作',
                dataIndex: 'operation',
                render: (text: any, record: any, index: number) => {
                    return (
                        <div className="editable-row-operations">
                            <AccountForm
                                title = "编辑"
                                data = {record}
                                onCreate = {values => this.onChange(values)}
                            />
                            <Button type="primary" onClick={() => this.onDelete(record)}>删除</Button>
                        </div>
                    );
                },
            }];
        return (
            <Table
                columns={columns}
                dataSource={record.data}
                pagination={false}
            />
        );
    }
    onChange(data: IAccountParam) {
        // same as source
        if (this.props.onChange) this.props.onChange(data);
        // console.log("change: ", values)
    }
    onDelete(data: IAccountParam) {
        // same as source
        if (this.props.onDelete) this.props.onDelete(data);
        // console.log("delete: ", values)
    }
    onCreate(data: IAccountParam) {
        // same as source
        if (this.props.onCreate) this.props.onCreate(data);
        // console.log("create: ", values)
    }
    render() {
        return (
            <div>
                <div style={{ marginBottom: 16 }}>
                    <Row gutter={16} align="bottom">
                        <Col span={2}>
                            <AccountForm
                                title = "账户"
                                onCreate = {values => this.onCreate(values)}
                            />
                        </Col>
                        <Col span={2}>
                            <Icon type="heart" style={{ color: '#FF0040' }} /> {this.state.amount}
                        </Col>
                    </Row>
                </div>
                <Table
                    columns={this.columns}
                    dataSource={this.state.data}
                    expandedRowRender={(record) => this.expandedRowRender(record)}
                    pagination={false}
                />
            </div>
        );
    }
}

export default AccountTable;