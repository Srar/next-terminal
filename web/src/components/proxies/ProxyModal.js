import React, {useState} from 'react';
import {Form, Input, Modal, Select, InputNumber} from "antd/lib/index";
import {isEmpty} from "../../utils/utils";
import ProxyTypes from './ProxyTypes';

const ProxyModal = ({title, visible, handleOk, handleCancel, confirmLoading, model}) => {

    const [form] = Form.useForm();

    const formItemLayout = {
        labelCol: {span: 6},
        wrapperCol: {span: 14},
    };

    if (model === null || model === undefined) {
        model = {}
    }

    if (isEmpty(model.type)) {
        model.type = 'socks5';
    }

    for (let key in model) {
        if (model.hasOwnProperty(key)) {
            if (model[key] === '-') {
                model[key] = '';
            }
        }
    }

    let [,setType] = useState(model.type);

    const handleAccountTypeChange = v => {
        setType(v);
        model.type = v;
    }

    return (

        <Modal
            title={title}
            visible={visible}
            maskClosable={false}

            onOk={() => {
                form
                    .validateFields()
                    .then(async (values) => {
                        const result = await handleOk(values);
                        if (result) {
                            form.resetFields();
                        }
                    })
                    .catch(info => {});
            }}
            onCancel={handleCancel}
            confirmLoading={confirmLoading}
            okText='确定'
            cancelText='取消'
        >

            <Form form={form} {...formItemLayout} initialValues={model}>
                <Form.Item name='id' noStyle>
                    <Input hidden={true}/>
                </Form.Item>

                <Form.Item label="代理名称" name='name' rules={[{required: true, message: '请输入代理名称'}]}>
                    <Input/>
                </Form.Item>

                <Form.Item label="代理类型" name='type' rules={[{required: true, message: '请选择代理类型'}]}>
                    <Select onChange={handleAccountTypeChange}>
                        {ProxyTypes.map(item => {
                            if (item.text === "无") 
                                return null;
                            return (<Select.Option key={item.value} value={item.value}>{item.text}</Select.Option>)
                        })}
                    </Select>
                </Form.Item>

                <Form.Item label="代理IP" name='host' rules={[{required: true, message: '请输入代理IP'}]}>
                    <Input/>
                </Form.Item>
                <Form.Item label="代理端口" name='port' rules={[{required: true, message: '请输入代理端口'}]}>
                    <InputNumber min={1} max={65535}/>
                </Form.Item>
                <Form.Item label="代理认证用户名" name='username'>
                    <Input/>
                </Form.Item>
                <Form.Item label="代理认证密码" name='password'>
                    <Input.Password/>
                </Form.Item>

            </Form>
        </Modal>
    )
};

export default ProxyModal;
