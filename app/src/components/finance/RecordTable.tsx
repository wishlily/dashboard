import React, {Component} from 'react';
import {
    Table,
    Button,
    DatePicker,
    Row,
    Col
} from 'antd';
import moment from 'moment';
import {IRecordParam} from '../../axios';
import RecordForm, {IRecordForm} from './RecordForm';

import {
    getFinanceRecord,
    setFinanceRecord,
    getFinanceAccount
} from '../../axios';

const { MonthPicker } = DatePicker;
const timeFmt = 'YYYY-MM-DD HH:mm:ss';

type RecordTableProps = {
}
interface RecordTableState {
    count: number;
    time: string;
    data: Array<IRecordParam>;
    account: Array<string>;
}

class RecordTable extends Component<RecordTableProps, RecordTableState> {
    constructor(props: any) {
        super(props);
        this.state = {
            count: 0,
            time: moment().format('YYYY-MM'),
            data: [],
            account: [],
        }
        this.columns = [{
            title: '交易类型',
            dataIndex: 'type',
            width: '10%',
        }, {
            title: '日期',
            dataIndex: 'time',
            width: '15%',
        }, {
            title: '金额',
            dataIndex: 'amount',
            width: '10%',
        }, {
            title: '备注',
            dataIndex: 'note',
            width: '45%',
        }, {
            title: '操作',
            dataIndex: 'operation',
            render: (text: any, record: any, index: number) => {
                return (
                    <div className="editable-row-operations">
                        <RecordForm
                            title = "编辑"
                            data = {this.state.data[index]}
                            account = {this.state.account}
                            onCreate = {values => this.onChange(values, index)}
                        />
                        <Button type="primary" onClick={() => this.onDelete(index)}>删除</Button>
                    </div>
                );
            },
        }];
    }
    componentDidMount() {
        getFinanceAccount("list").then(data =>{
            this.setState({account: data})
        })
        this.getData(this.state.time)
    }
    componentDidUpdate(prevProps: RecordTableProps, prevState: RecordTableState) {
        if (this.state.count !== prevState.count || this.state.time !== prevState.time) {
            this.getData(prevState.time)
            // console.log("UPDATE!!!")
        }
    }
    columns: any;
    getData(tmonth: string) {
        const time = moment(tmonth + '-01 00:00:00', timeFmt)
        const t_start = time.format(timeFmt)
        const t_end = time.add(1, 'months').subtract(1, 'seconds').format(timeFmt)
        // console.log("start: ", t_start, ", end: ", t_end)
        getFinanceRecord(t_start, t_end).then(data => {
            this.setState({ data: data })
            // console.log("get: ", data)
        });
    }
    onChange(values: IRecordForm, index: number) {
        const datas = this.state.data
        const data = datas[index]
        data.type = values.type
        data.time = values.time.format(timeFmt)
        data.account[0] = values.accountM
        if (values.accountS !== undefined ) data.account[1] = values.accountS
        data.amount = values.amount
        if (values.member !== undefined ) data.member = values.member
        if (values.class !== undefined ) data.class = values.class
        data.note = values.note
        datas[index] = data
        // this.setState({ data: datas })
        // console.log("change: ", data)
        setFinanceRecord('chg', data).then(res =>{
            if (res && res.message === 'ok') {
                this.setState({count: this.state.count+1})
            }
        })
    }
    onDelete(index: number) {
        const datas = this.state.data
        const data = datas[index]
        // console.log("delete: ", data);
        setFinanceRecord('del', data).then(res =>{
            if (res && res.message === 'ok') {
                this.setState({count: this.state.count+1})
            }
        })
    }
    onCreate(values: IRecordForm) {
        const data: IRecordParam = {
            key: '',
            type: values.type,
            time: values.time.format(timeFmt),
            account: [values.accountM],
            amount: values.amount,
            member: values.member,
            class: values.class,
            note: values.note
        }
        if (values.accountS !== undefined) data.account[1] = values.accountS
        // console.log("create: ", data)
        setFinanceRecord('add', data).then(res =>{
            if (res && res.message === 'ok') {
                this.setState({count: this.state.count+1})
            }
        })
    }
    onChangeTime(date: any, dateString: string) {
        this.setState({ time: dateString });
        // console.log(dateString);
    }
    render() {
        return (
            <div>
                <div style={{ marginBottom: 16 }}>
                    <Row gutter={16}>
                        <Col span={2}>
                            <MonthPicker
                                defaultValue={moment(this.state.time, 'YYYY-MM')}
                                onChange={(date, dateString) => this.onChangeTime(date, dateString)}
                                placeholder="请选择月份"
                                style={{ width: 100 }}
                            />
                        </Col>
                        <Col span={2}>
                            <RecordForm
                                title = "新建"
                                account = {this.state.account}
                                onCreate = {values => this.onCreate(values)}
                            />
                        </Col>
                    </Row>
                </div>
                <Table bordered dataSource={this.state.data} columns={this.columns} />
            </div>
        );
    }
}

export default RecordTable;