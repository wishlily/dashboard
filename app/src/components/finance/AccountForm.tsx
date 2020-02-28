import React, {Component} from 'react';
import {
    Button,
    Modal,
    Form,
    DatePicker,
    InputNumber,
    Input,
    Select
} from 'antd';
import moment from 'moment';
import { FormComponentProps } from 'antd/lib/form/Form';
import { IAccountParam, minVaildTS } from '../../axios';
import { SelectValue } from 'antd/lib/select';

interface IAccountForm {
    account: string; // id
    type: string;
    class: string;
    amount: number;
    unit: number;
    nuv: number;
    deadline: moment.Moment;
}

const FormItem = Form.Item;
const Option = Select.Option;

const timeFmt = 'YYYY-MM-DD HH:mm:ss';
const classData = ['零钱', '活期', '信贷', '基金', '定期'];
const typeData = ['CNY', 'USD'];

type BasicFormsProps = {
    visible: boolean;
    onCancel: () => void;
    onCreate: () => void;
    data: IAccountParam;
} & FormComponentProps;
interface BasicFormsState {
    data: IAccountParam
}

class BasicForms extends Component<BasicFormsProps, BasicFormsState> {
    state = {
        data: {
            time: moment().format(timeFmt),
            id: '',
            type: typeData[0],
            amount: 0,
            unit: 0,
            nuv: 1,
            class: classData[0],
            deadline: moment(0).format(timeFmt)
        },
    };
    componentDidUpdate(prevProps: BasicFormsProps) {
        if (prevProps.data && prevProps.data !== this.state.data) {
            this.setState({ data: prevProps.data });
            let deadine: moment.Moment = moment(0)
            if (prevProps.data.deadline) {
                deadine = moment(prevProps.data.deadline, timeFmt)
            }
            prevProps.form.setFieldsValue({
                type: prevProps.data.type,
                deadline: deadine,
                account: prevProps.data.id,
                class: prevProps.data.class,
                unit: prevProps.data.unit,
                nuv: prevProps.data.nuv,
                amount: prevProps.data.amount,
            });
        }
    }
    classChange(key: SelectValue) {
        const data = this.state.data
        data.class = key.toString()
        this.setState({ data: data })
    }
    typeChange(key: SelectValue) {
        const data = this.state.data
        data.type = key.toString()
        this.setState({ data: data })
    }
    render() {
        const { visible, onCancel, onCreate, form } = this.props;
        const { getFieldDecorator } = form;

        const classOptions = classData.map(val => <Option key={val}>{val}</Option>);
        const typeOptions = typeData.map(val => <Option key={val}>{val}</Option>);
        const formItemLayout = {
            labelCol: { span: 4 },
            wrapperCol: { span: 16 },
        };
        return (
            <Modal
                title="账户信息"
                visible={visible}
                okText="完成"
                cancelText="取消"
                onCancel={onCancel}
                onOk={onCreate}
            >
                <Form>
                    <FormItem {...formItemLayout} label="类型" colon={false}>
                        {getFieldDecorator('class', {
                            initialValue: this.state.data.class,
                            rules: [{required: true, message: '请输入账户类型'}]
                        })(
                            <Select style={{ width: '100%' }} onChange={ (v) => this.classChange(v) }>
                                {classOptions}
                            </Select>
                        )}
                    </FormItem>
                    <FormItem {...formItemLayout} label="账户" colon={false}>
                        {getFieldDecorator('account', {
                            initialValue: this.state.data.id,
                            rules: [{required: true, message: '请输入账户名'}],
                        })(
                            <Input style={{ width: '100%' }} />
                        )}
                    </FormItem>
                    <FormItem {...formItemLayout} label="币种" colon={false}>
                        {getFieldDecorator('type', { initialValue: this.state.data.type })(
                            <Select style={{ width: '100%' }} onChange={ (v) => this.typeChange(v) }>
                                {typeOptions}
                            </Select>
                        )}
                    </FormItem>
                    <FormItem {...formItemLayout} label="份额" colon={false}>
                        {getFieldDecorator('unit', { initialValue: this.state.data.unit })(
                            <InputNumber style={{ width: '100%' }} min={0} step={0.01} />
                        )}
                    </FormItem>
                    <FormItem {...formItemLayout} label="净值" colon={false}>
                        {getFieldDecorator('nuv', { initialValue: this.state.data.nuv })(
                            <InputNumber style={{ width: '100%' }} min={0} step={0.01} />
                        )}
                    </FormItem>
                    <FormItem {...formItemLayout} label="投入" colon={false}>
                        {getFieldDecorator('amount', { initialValue: this.state.data.amount })(
                            <InputNumber style={{ width: '100%' }} min={0} step={0.01} />
                        )}
                    </FormItem>
                    <FormItem {...formItemLayout} label="截至日期" colon={false}>
                        {getFieldDecorator('deadline', { initialValue: moment(this.state.data.deadline) })(
                            <DatePicker style={{ width: '100%' }} showTime format={timeFmt} />
                        )}
                    </FormItem>
                </Form>
            </Modal>
        )
    }
}

const CollectionCreateForm: any = Form.create()(BasicForms);

type AccountFormProps = {
    title: string;
    onClick?: () => void;
    onCreate?: (v: IAccountParam) => void;
    data?: IAccountParam;
}

class AccountForm extends Component<AccountFormProps> {
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
        form.validateFields((err: any, values: IAccountForm) => {
            if (err) {
                return;
            }
            const data: IAccountParam = {
                time: moment().format(timeFmt),
                id: values.account,
                type: values.type,
                class: values.class,
                amount: values.amount
            }
            if (values.unit && values.unit > 0) {
                data.unit = values.unit
            }
            if (values.nuv && values.nuv > 0) {
                data.nuv = values.nuv
            }
            if (values.deadline) {
                let t = moment(values.deadline).unix()
                if (t > minVaildTS()) {
                    data.deadline = values.deadline.format(timeFmt)
                }
            }
            if (this.props.onCreate) this.props.onCreate(data);
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
                />
            </span>
        );
    }
}

export default AccountForm;