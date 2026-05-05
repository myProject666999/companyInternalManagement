import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, message, Space, Popconfirm, Card, Upload, Descriptions, Divider, Tag } from 'antd';
import { PlusOutlined, DeleteOutlined, DownloadOutlined, EyeOutlined, SearchOutlined, UploadOutlined } from '@ant-design/icons';
import api from '../utils/axios';
import { Document, User } from '../types';

const { Option } = Select;

const Documents: React.FC = () => {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailVisible, setDetailVisible] = useState(false);
  const [viewingDoc, setViewingDoc] = useState<Document | null>(null);
  const [searchForm] = Form.useForm();
  const [fileList, setFileList] = useState<any[]>([]);

  const fetchDocuments = async (params?: any) => {
    setLoading(true);
    try {
      const response = await api.get('/documents', { params });
      setDocuments(response.data);
    } catch (error) {
      message.error('获取文档列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDocuments();
  }, []);

  const handleSearch = (values: any) => {
    fetchDocuments(values);
  };

  const handleView = async (record: Document) => {
    try {
      const response = await api.get(`/documents/${record.id}`);
      setViewingDoc(response.data);
      setDetailVisible(true);
    } catch (error) {
      message.error('获取文档详情失败');
    }
  };

  const handleDownload = async (record: Document) => {
    try {
      const response = await api.get(`/documents/${record.id}/download`, {
        responseType: 'blob',
      });
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', record.file_path.split('/').pop() || record.title);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      message.success('下载成功');
    } catch (error) {
      message.error('下载失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/documents/${id}`);
      message.success('删除成功');
      fetchDocuments();
    } catch (error: any) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  const handleUploadChange = (info: any) => {
    setFileList(info.fileList);
  };

  const handleUpload = async () => {
    if (fileList.length === 0) {
      message.warning('请先选择文件');
      return;
    }

    const formData = new FormData();
    formData.append('file', fileList[0].originFileObj);

    try {
      await api.post('/documents', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      message.success('上传成功');
      setFileList([]);
      setModalVisible(false);
      fetchDocuments();
    } catch (error: any) {
      message.error(error.response?.data?.error || '上传失败');
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '文档标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
    {
      title: '文件大小',
      dataIndex: 'file_size',
      key: 'file_size',
      width: 120,
      render: (size: number) => formatFileSize(size),
    },
    {
      title: '上传者',
      dataIndex: ['uploader', 'name'],
      key: 'uploader',
      width: 100,
    },
    {
      title: '下载次数',
      dataIndex: 'download_count',
      key: 'download_count',
      width: 100,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width: 220,
      render: (_: any, record: Document) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleView(record)}
          >
            查看
          </Button>
          <Button
            type="link"
            icon={<DownloadOutlined />}
            onClick={() => handleDownload(record)}
          >
            下载
          </Button>
          <Popconfirm
            title="确定要删除这个文档吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Form
          form={searchForm}
          layout="inline"
          onFinish={handleSearch}
        >
          <Form.Item name="title" label="文档标题">
            <Input placeholder="请输入文档标题" style={{ width: 200 }} />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
                搜索
              </Button>
              <Button onClick={() => { searchForm.resetFields(); fetchDocuments(); }}>
                重置
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>

      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
          上传文档
        </Button>
      </div>
      
      <Table
        columns={columns}
        dataSource={documents}
        rowKey="id"
        loading={loading}
        bordered
      />

      <Modal
        title="上传文档"
        open={modalVisible}
        onCancel={() => { setModalVisible(false); setFileList([]); }}
        onOk={handleUpload}
        okText="上传"
        cancelText="取消"
        width={500}
      >
        <div style={{ padding: '20px 0' }}>
          <Upload
            fileList={fileList}
            onChange={handleUploadChange}
            beforeUpload={() => false}
            maxCount={1}
          >
            <Button icon={<UploadOutlined />}>选择文件</Button>
          </Upload>
          <div style={{ marginTop: 16, color: '#666' }}>
            <p>支持的文件格式：所有文档类型</p>
            <p>最大文件大小：50MB</p>
          </div>
        </div>
      </Modal>

      <Modal
        title="文档详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={[
          <Button key="download" type="primary" icon={<DownloadOutlined />} onClick={() => { if (viewingDoc) handleDownload(viewingDoc); }}>
            下载文档
          </Button>,
          <Button key="close" onClick={() => setDetailVisible(false)}>
            关闭
          </Button>,
        ]}
        width={600}
      >
        {viewingDoc && (
          <div>
            <Descriptions bordered size="small" column={2}>
              <Descriptions.Item label="文档标题" span={2}>
                {viewingDoc.title}
              </Descriptions.Item>
              <Descriptions.Item label="文件大小">
                {formatFileSize(viewingDoc.file_size)}
              </Descriptions.Item>
              <Descriptions.Item label="下载次数">
                {viewingDoc.download_count}
              </Descriptions.Item>
              <Descriptions.Item label="上传者">
                {viewingDoc.uploader?.name}
              </Descriptions.Item>
              <Descriptions.Item label="上传时间">
                {viewingDoc.created_at}
              </Descriptions.Item>
            </Descriptions>
            <Divider />
            <div style={{ textAlign: 'center', padding: '40px 0' }}>
              <p style={{ color: '#666' }}>文档预览功能正在开发中</p>
              <p style={{ color: '#999', fontSize: 12 }}>请点击"下载文档"按钮查看完整内容</p>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default Documents;
