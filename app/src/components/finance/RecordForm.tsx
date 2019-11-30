import React, {Component} from 'react';
import {
    Button,
    Modal,
    Form,
    DatePicker,
    Tabs,
    InputNumber,
    Cascader,
    Input,
    Select
} from 'antd';
import moment from 'moment';
import {FormComponentProps} from 'antd/lib/form/Form';
import {IRecordParam} from '../../axios';

interface IClass {
    [key: string]: any
}

export interface IRecordForm {
    type: string;
    class?: Array<string>,
    time: moment.Moment,
    accountM: string,
    accountS?: string,
    amount: number,
    member?: string,
    note: string,
}

const FormItem = Form.Item;
const TabPane = Tabs.TabPane;
const { TextArea } = Input;
const { Option } = Select;

const timeFmt = 'YYYY-MM-DD HH:mm:ss';
const typeData = {TYPEI: '支出', TYPEO: '收入', TYPER: '转账', TYPEL: '借出', TYPEB: '借入'};
const {TYPEI, TYPEO, TYPER, TYPEL, TYPEB} = typeData;
const classData: IClass = require("./classify.json");

const classOptions = Object.keys(classData).map(v => ({
    value: v,
    label: v,
    children: Object.keys(classData[v]).map(v1 => ({
        value: v1,
        label: v1,
        children: (classData[v][v1] as Array<string>).map(v2 =>({
            value: v2,
            label: v2,
        }))
    }))
}));

type BasicFormsProps = {
    visible: boolean;
    onCancel: () => void;
    onCreate: () => void;
    account: Array<string>;
    data: IRecordParam;
} & FormComponentProps;
interface BasicFormsState {
    data: IRecordParam;
}

class BasicForms extends Component<BasicFormsProps, BasicFormsState> {
    constructor(props: any) {
        super(props);
        const m = Object.keys(classData[TYPEI])[0];
        const s = classData[TYPEI][m][0]
        this.state = {
            data: {
                key: '',
                type: TYPEI,
                amount: 0,
                account: [],
                time: moment().format(timeFmt),
                class: [m, s],
            }
        }
    }
    componentDidUpdate(prevProps: BasicFormsProps) {
        if (prevProps.data && prevProps.data !== this.state.data) {
            this.setState({ data: prevProps.data });
            prevProps.form.setFieldsValue({
                type: prevProps.data.type,
                class: prevProps.data.class,
                time: moment(prevProps.data.time, timeFmt),
                accountM: prevProps.data.account[0],
                accountS: prevProps.data.account[1],
                amount: prevProps.data.amount,
                member: prevProps.data.member,
                note: prevProps.data.note,
            });
        }
    }
    onTabsChange = (key: string) => {
        if (key !== TYPEI && key !== TYPEO) return;
        const m = Object.keys(classData[key])[0];
        const s = classData[key][m][0]
        this.setState({
            data: Object.assign({}, this.state.data, {
                type: key,
                class: [m, s],
            })
        });
        this.props.form.setFieldsValue({
            class: [m, s],
        });
    };
    render() {
        const { visible, onCancel, onCreate, form, account } = this.props;
        const { getFieldDecorator } = form;

        var accountData: Array<string> = [];
        if (account instanceof Array) {
            accountData = account
        }
        const accountOptions = accountData.map(val => <Option key={val}>{val}</Option>);
        const formItemLayout = {
            labelCol: { span: 4 },
            wrapperCol: { span: 16 },
        };
        return (
            <Modal
                visible={visible}
                okText="完成"
                cancelText="取消"
                onCancel={onCancel}
                onOk={onCreate}
            >
                <FormItem>
                    {getFieldDecorator('type', { initialValue: this.state.data.type })(
                        <Tabs defaultActiveKey={this.state.data.type} onChange={ this.onTabsChange } type="card">
                            <TabPane tab={ TYPEI } key={ TYPEI }>
                                <Form>
                                    <FormItem {...formItemLayout} label="金额" colon={false}>
                                        {getFieldDecorator('amount', { initialValue: this.state.data.amount })(
                                            <InputNumber style={{ width: '100%' }} min={0} step={0.01} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="时间" colon={false}>
                                        {getFieldDecorator('time', { initialValue: moment(this.state.data.time, timeFmt) })(
                                            <DatePicker style={{ width: '100%' }} showTime format={timeFmt} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="分类" colon={false}>
                                        {getFieldDecorator('class', { initialValue: this.state.data.class })(
                                            <Cascader style={{ width: '100%' }} options={classOptions[0].children} changeOnSelect />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="账户" colon={false}>
                                        {getFieldDecorator('accountM', { initialValue: this.state.data.account[0] })(
                                            <Select showSearch style={{ width: '100%' }} >
                                                {accountOptions}
                                            </Select>
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="备注" colon={false}>
                                        {getFieldDecorator('note', { initialValue: this.state.data.note })(
                                            <TextArea style={{ width: '100%' }} autosize />
                                        )}
                                    </FormItem>
                                </Form>
                            </TabPane>
                            <TabPane tab={ TYPEO } key={ TYPEO }>
                                <Form>
                                    <FormItem {...formItemLayout} label="金额" colon={false}>
                                        {getFieldDecorator('amount', { initialValue: this.state.data.amount })(
                                            <InputNumber style={{ width: '100%' }} min={0} step={0.01} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="时间" colon={false}>
                                        {getFieldDecorator('time', { initialValue: moment(this.state.data.time, timeFmt) })(
                                            <DatePicker style={{ width: '100%' }} showTime format={timeFmt} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="分类" colon={false}>
                                        {getFieldDecorator('class', { initialValue: this.state.data.class })(
                                            <Cascader style={{ width: '100%' }} options={classOptions[1].children} changeOnSelect />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="账户" colon={false}>
                                        {getFieldDecorator('accountM', { initialValue: this.state.data.account[0] })(
                                            <Select showSearch style={{ width: '100%' }} >
                                                {accountOptions}
                                            </Select>
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="备注" colon={false}>
                                        {getFieldDecorator('note', { initialValue: this.state.data.note })(
                                            <TextArea style={{ width: '100%' }} autosize />
                                        )}
                                    </FormItem>
                                </Form>
                            </TabPane>
                            <TabPane tab={ TYPER } key={ TYPER }>
                                <Form>
                                    <FormItem {...formItemLayout} label="金额" colon={false}>
                                        {getFieldDecorator('amount', { initialValue: this.state.data.amount })(
                                            <InputNumber style={{ width: '100%' }} min={0} step={0.01} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="时间" colon={false}>
                                        {getFieldDecorator('time', { initialValue: moment(this.state.data.time, timeFmt) })(
                                            <DatePicker style={{ width: '100%' }} showTime format={timeFmt} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="转出" colon={false}>
                                        {getFieldDecorator('accountM', { initialValue: this.state.data.account[0] })(
                                            <Select showSearch style={{ width: '100%' }} >
                                                {accountOptions}
                                            </Select>
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="转入" colon={false}>
                                        {getFieldDecorator('accountS', { initialValue: this.state.data.account[1] })(
                                            <Select showSearch style={{ width: '100%' }} >
                                                {accountOptions}
                                            </Select>
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="备注" colon={false}>
                                        {getFieldDecorator('note', { initialValue: this.state.data.note })(
                                            <TextArea
                                                style={{ width: '100%' }}
                                                placeholder="#份额#xx"
                                                autosize
                                            />
                                        )}
                                    </FormItem>
                                </Form>
                            </TabPane>
                            <TabPane tab={ TYPEL } key={ TYPEL }>
                                <Form>
                                    <FormItem {...formItemLayout} label="金额" colon={false}>
                                        {getFieldDecorator('amount', { initialValue: this.state.data.amount })(
                                            <InputNumber style={{ width: '100%' }} min={0} step={0.01} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="时间" colon={false}>
                                        {getFieldDecorator('time', { initialValue: moment(this.state.data.time, timeFmt) })(
                                            <DatePicker style={{ width: '100%' }} showTime format={timeFmt} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="借贷人" colon={false}>
                                        {getFieldDecorator('member', { initialValue: this.state.data.member })(
                                            <Input style={{ width: '100%' }} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="账户" colon={false}>
                                        {getFieldDecorator('accountM', { initialValue: this.state.data.account[0] })(
                                            <Select showSearch style={{ width: '100%' }} >
                                                {accountOptions}
                                            </Select>
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="备注" colon={false}>
                                        {getFieldDecorator('note', { initialValue: this.state.data.note })(
                                            <TextArea
                                                style={{ width: '100%' }}
                                                placeholder="#还款时间#xxxx-xx-xx xx:xx:xx #还款金额#xx"
                                                autosize
                                            />
                                        )}
                                    </FormItem>
                                </Form>
                            </TabPane>
                            <TabPane tab={ TYPEB } key={ TYPEB }>
                                <Form>
                                    <FormItem {...formItemLayout} label="金额" colon={false}>
                                        {getFieldDecorator('amount', { initialValue: this.state.data.amount })(
                                            <InputNumber style={{ width: '100%' }} min={0} step={0.01} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="时间" colon={false}>
                                        {getFieldDecorator('time', { initialValue: moment(this.state.data.time, timeFmt) })(
                                            <DatePicker style={{ width: '100%' }} showTime format={timeFmt} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="借贷人" colon={false}>
                                        {getFieldDecorator('member', { initialValue: this.state.data.member })(
                                            <Input style={{ width: '100%' }} />
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="账户" colon={false}>
                                        {getFieldDecorator('accountM', { initialValue: this.state.data.account[0] })(
                                            <Select showSearch style={{ width: '100%' }} >
                                                {accountOptions}
                                            </Select>
                                        )}
                                    </FormItem>
                                    <FormItem {...formItemLayout} label="备注" colon={false}>
                                        {getFieldDecorator('note', { initialValue: this.state.data.note })(
                                            <TextArea
                                                style={{ width: '100%' }}
                                                placeholder="#还款时间#xxxx-xx-xx xx:xx:xx #还款金额#xx"
                                                autosize
                                            />
                                        )}
                                    </FormItem>
                                </Form>
                            </TabPane>
                        </Tabs>
                    )}
                </FormItem>
            </Modal>
        )
    }
}

const CollectionCreateForm: any = Form.create()(BasicForms);

type RecordFormProps = {
    title: string;
    onClick?: () => void;
    onCreate?: (v: IRecordForm) => void;
    account: Array<string>;
    data?: IRecordParam;
}

class RecordForm extends Component<RecordFormProps> {
    state = {
        visible: false,
    };
    form: any;
    showModal = () => {
        if (this.props.onClick) this.props.onClick();
        this.setState({ visible: true });
    };
    handleCancel = () => {
        this.setState({ visible: false });
    };
    handleCreate = () => {
        const form = this.form;
        form.validateFields((err: any, values: IRecordForm) => {
            if (err) {
                return;
            }
            // console.log('Received values of form: ', values);
            if (this.props.onCreate) this.props.onCreate(values);
            form.resetFields();
            this.setState({ visible: false });
        });
    };
    saveFormRef = (form: any) => {
        this.form = form;
    };
    render() {
        return (
            <span>
                <Button type="primary" onClick={this.showModal}>{this.props.title}</Button>
                <CollectionCreateForm
                    ref={this.saveFormRef}
                    visible={this.state.visible}
                    onCancel={this.handleCancel}
                    onCreate={this.handleCreate}
                    data={this.props.data}
                    account={this.props.account}
                />
            </span>
        );
    }
}

export default RecordForm;