import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate, Link } from 'react-router-dom';
import { Layout, Menu, Button, message, Spin, Card, Row, Col, Form, Input, InputNumber, Upload, Modal, Image, Select } from 'antd';
import { UploadOutlined, ShoppingCartOutlined, PlayCircleOutlined, LoginOutlined, LogoutOutlined, VideoCameraOutlined, HomeOutlined, AppstoreOutlined, UserOutlined } from '@ant-design/icons';
import ReactPlayer from 'react-player';
import axios from 'axios';

// 布局组件
const { Header, Content, Footer, Sider } = Layout;
const { Option } = Select;

// 全局Axios配置
axios.defaults.baseURL = '/api';
axios.defaults.timeout = 60000;

axios.interceptors.request.use(config => {
  const token = localStorage.getItem('geekedu_token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
}, err => {
  message.error('请求发送失败，请检查网络');
  return Promise.reject(err);
});

axios.interceptors.response.use(
  res => res,
  err => {
    const errMsg = err.response?.data?.msg || err.message || '网络异常，请稍后重试';
    message.error(errMsg);
    return Promise.reject(err);
  }
);

const App = () => {
  // 全局状态
  const [userInfo, setUserInfo] = useState(null);
  const [loading, setLoading] = useState(false);
  const [courses, setCourses] = useState([]);
  const [loginVisible, setLoginVisible] = useState(false);
  const [isLoginTab, setIsLoginTab] = useState(true);
  const [form] = Form.useForm();
  const [registerForm] = Form.useForm();
  const [coverFile, setCoverFile] = useState(null);
  const [coverPreview, setCoverPreview] = useState('');
  const [courseForm] = Form.useForm();
  const [videoFile, setVideoFile] = useState(null);
  const [videoForm] = Form.useForm();
  const [videoModalVisible, setVideoModalVisible] = useState(false);
  const [courseVideos, setCourseVideos] = useState([]);
  const [currentPlayUrl, setCurrentPlayUrl] = useState('');
  const [currentVideoTitle, setCurrentVideoTitle] = useState('');
  const [playerLoading, setPlayerLoading] = useState(false);
  const [currentCourseId, setCurrentCourseId] = useState(null);
  const [collapsed, setCollapsed] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('geekedu_token');
    const uid = localStorage.getItem('geekedu_uid');
    const role = localStorage.getItem('geekedu_role');
    if (token && uid && role) {
      setUserInfo({ 
        token, 
        uid: parseInt(uid, 10),
        role: parseInt(role, 10) 
      });
    }
    fetchCourses();
  }, []);

  const fetchCourses = async () => {
    setLoading(true);
    try {
      const res = await axios.get('/v1/courses');
      if (res?.data?.code === 200) {
        setCourses(res?.data?.data?.courses || []);
      }
    } catch (err) {
      console.error('获取课程列表失败：', err);
      setCourses([]);
    } finally {
      setLoading(false);
    }
  };

  const handleLogin = async (values) => {
    try {
      message.loading('登录中...', 0.5);
      const res = await axios.post('/v1/auth/login', values);
      
      if (!res || !res.data) {
        message.error('登录响应数据异常，请重试');
        return;
      }
      
      if (res?.data?.code !== 200) {
        message.error(res?.data?.msg || '登录失败，请检查用户名或密码');
        return;
      }
      
      const loginData = res?.data?.data || {};
      const { uid, token, role } = loginData;
  
      if (uid === undefined || uid === null || 
          token === undefined || token === null || 
          role === undefined || role === null) {
        message.error('登录数据不完整，请重试（具体：uid/token/role缺失）');
        return;
      }
      
      const uidStr = uid.toString();
      const roleStr = role.toString();
      
      localStorage.setItem('geekedu_token', token);
      localStorage.setItem('geekedu_uid', uidStr);
      localStorage.setItem('geekedu_role', roleStr);
      
      setUserInfo({ 
        token, 
        uid: parseInt(uidStr, 10), 
        role: parseInt(roleStr, 10) 
      });
      setLoginVisible(false);
      form.resetFields();
      registerForm.resetFields();
      message.success('登录成功！');
      fetchCourses();
    } catch (err) {
      console.error('登录失败：', err);
      message.error('登录失败，请检查用户名或密码是否正确');
      form.resetFields();
    }
  };

  const handleRegister = async (values) => {
    try {
      message.loading('注册中...', 0.5);
      const res = await axios.post('/v1/auth/register', {
        username: values.username,
        password: values.password,
      });
      if (res?.data?.code === 200) {
        message.success(res?.data?.msg || '注册成功，请切换到登录页登录');
        setIsLoginTab(true);
        registerForm.resetFields();
      }
    } catch (err) {
      console.error('注册失败：', err);
      message.error('注册失败，用户名可能已存在');
      registerForm.resetFields();
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('geekedu_token');
    localStorage.removeItem('geekedu_uid');
    localStorage.removeItem('geekedu_role');
    setUserInfo(null);
    setCourses([]);
    setVideoModalVisible(false);
    setCourseVideos([]);
    setCurrentPlayUrl('');
    message.success('退出登录成功');
  };

  const handleBuyCourse = async (courseId) => {
    if (!userInfo) {
      message.warning('请先登录再购买课程');
      setLoginVisible(true);
      return;
    }
    try {
      const res = await axios.post('/v1/orders', { course_id: courseId });
      if (res?.data?.code === 200) {
        message.success('购买成功，可播放该课程视频');
        fetchCourses();
      }
    } catch (err) {
      console.error('购买课程失败：', err);
      message.error('购买课程失败，请稍后重试');
    }
  };

  const getCourseVideos = async (courseId) => {
    if (!courseId || !userInfo) {
      return;
    }
    setPlayerLoading(true);
    try {
      const res = await axios.get(`/v1/courses/${courseId}/videos`);
      if (res?.data?.code === 200) {
        setCourseVideos(res?.data?.data?.videos || []);
        setCurrentCourseId(courseId);
        setVideoModalVisible(true);
        const firstVideo = res?.data?.data?.videos?.[0];
        if (firstVideo) {
          handleGetVideoPlayUrl(firstVideo.id, firstVideo.title);
        }
      }
    } catch (err) {
      console.error('获取课程视频列表失败：', err);
      message.error('获取视频列表失败，可能未购买该课程');
    } finally {
      setPlayerLoading(false);
    }
  };

  const handleGetVideoPlayUrl = async (videoId, videoTitle) => {
    if (!videoId || !userInfo) {
      return;
    }
    setPlayerLoading(true);
    try {
      const res = await axios.get(`/v1/player/${videoId}`);
      if (res?.data?.code === 200) {
        setCurrentPlayUrl(res?.data?.data?.signed_url || '');
        setCurrentVideoTitle(videoTitle || '课程视频');
      }
    } catch (err) {
      console.error('获取视频播放地址失败：', err);
      message.error('获取播放地址失败，请稍后重试');
    } finally {
      setPlayerLoading(false);
    }
  };

  const handlePlayCourse = async (courseId) => {
    if (!userInfo) {
      message.warning('请先登录再播放课程');
      setLoginVisible(true);
      return;
    }
    await getCourseVideos(courseId);
  };

  const renderVideoModalContent = () => {
    return (
      <div className="video-modal-content" style={{ display: 'flex', gap: '24px', height: '70vh' }}>
        <div className="video-list-container">
          <h3 className="video-list-title">视频列表</h3>
          {courseVideos.length > 0 ? (
            <div className="video-list-items">
              {courseVideos.map(video => (
                <div
                  key={video.id}
                  className={`video-list-item ${currentVideoTitle === video.title ? 'active' : ''}`}
                  onClick={() => handleGetVideoPlayUrl(video.id, video.title)}
                >
                  <p className="video-list-item-title">
                    {video.title}
                  </p>
                  <p className="video-list-item-time">
                    创建时间：{video.createdAt || '未知时间'}
                  </p>
                </div>
              ))}
            </div>
          ) : (
            <div className="video-list-empty">
              暂无该课程的视频数据
            </div>
          )}
        </div>

        <div className="video-player-container">
          {playerLoading ? (
            <div className="player-loading">
              <Spin size="large" tip="正在加载播放地址..." />
            </div>
          ) : currentPlayUrl ? (
            <div className="player-wrapper">
              <h3 className="player-title">{currentVideoTitle}</h3>
              <ReactPlayer
                url={currentPlayUrl}
                width="100%"
                height="calc(100% - 40px)"
                controls={true}
                playing={false}
                playbackRate={[0.5, 0.75, 1, 1.25, 1.5, 2]}
                volume={0.7}
                className="video-player"
              />
            </div>
          ) : (
            <div className="player-placeholder">
              <p className="placeholder-text">请选择左侧视频进行播放</p>
              <PlayCircleOutlined style={{ fontSize: 48, color: '#999', marginBottom: 16 }} />
            </div>
          )}
        </div>
      </div>
    );
  };

  const handleUploadCover = async () => {
    if (!coverFile || !userInfo || userInfo.role !== 1) {
      message.warning('请选择封面图片，且仅管理员可操作');
      return;
    }
    const formData = new FormData();
    formData.append('cover', coverFile);

    try {
      const res = await axios.post('/v1/courses/upload/cover', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      if (res?.data?.code === 200) {
        const coverOssKey = res?.data?.data?.cover_oss_key || '';
        message.success('封面上传成功');
        courseForm.setFieldsValue({ cover_oss_key: coverOssKey });
        setCoverPreview(URL.createObjectURL(coverFile));
      }
    } catch (err) {
      console.error('封面上传失败：', err);
      message.error('封面上传失败，请检查OSS配置');
    }
  };

  const handlePublishCourse = async (values) => {
    if (!userInfo || userInfo.role !== 1) {
      message.warning('仅管理员可发布课程');
      return;
    }

    if (!values.cover_oss_key) {
      message.warning('请先上传课程封面');
      return;
    }

    if (!values.title || values.price <= 0 || !values.intro) {
      message.warning('请填写完整且有效的课程信息');
      return;
    }

    try {
      const res = await axios.post('/v1/courses', values);
      if (res?.data?.code === 200) {
        message.success('课程发布成功');
        courseForm.resetFields();
        setCoverFile(null);
        setCoverPreview('');
        fetchCourses();
      }
    } catch (err) {
      console.error('课程发布失败：', err);
      message.error('课程发布失败，请稍后重试');
    }
  };

  const handleCoverChange = (info) => {
    const file = info.file;
    if (file) {
      setCoverFile(file);
      setCoverPreview(URL.createObjectURL(file));
    }
  };

  const handleVideoChange = (info) => {
    const file = info.file;
    if (file) {
      setVideoFile(file);
      message.info(`已选择视频：${file.name}`);
    }
  };

  const handleUploadVideo = async (values) => {
    if (!userInfo || userInfo.role !== 1) {
      message.warning('仅管理员可上传视频');
      return;
    }
    if (!videoFile) {
      message.warning('请先选择要上传的视频文件');
      return;
    }
    if (!values.course_id || !values.video_title) {
      message.warning('请选择课程并填写视频标题');
      return;
    }

    const formData = new FormData();
    formData.append('course_id', values.course_id);
    formData.append('title', values.video_title);
    formData.append('video', videoFile);

    try {
      const res = await axios.post('/v1/courses/upload/video', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
        onUploadProgress: (progressEvent) => {
          const progress = Math.round((progressEvent.loaded / progressEvent.total) * 100);
          message.info(`视频上传中：${progress}%`, 0.1);
        },
      });

      if (res?.data?.code === 200) {
        message.success('视频上传成功，已绑定对应课程');
        videoForm.resetFields();
        setVideoFile(null);
        fetchCourses();
      }
    } catch (err) {
      console.error('视频上传失败：', err);
      message.error('视频上传失败，请检查OSS配置或视频格式');
    }
  };

  return (
    <BrowserRouter>
      <Layout className="app-layout">
        <Header className="app-header">
          <div className="header-content">
            <div className="logo-section">
              <div className="logo">GE</div>
              <h1 className="app-title">在线视频学习平台</h1>
            </div>
            
            <div className="user-section">
              {userInfo ? (
                <div className="user-info">
                  <span className="user-avatar">
                    <UserOutlined />
                  </span>
                  <span className="user-name">
                    {userInfo.role === 1 ? '管理员' : '学员'} #{userInfo.uid}
                  </span>
                  <Button 
                    type="text" 
                    icon={<LogoutOutlined />}
                    onClick={handleLogout}
                    className="logout-btn"
                  >
                    退出
                  </Button>
                </div>
              ) : (
                <Button 
                  type="primary" 
                  icon={<LoginOutlined />}
                  onClick={() => {
                    setLoginVisible(true);
                    setIsLoginTab(true);
                  }}
                  className="login-btn"
                >
                  登录/注册
                </Button>
              )}
            </div>
          </div>
        </Header>

        <Layout className="main-layout">
          <Sider 
            collapsible 
            collapsed={collapsed} 
            onCollapse={setCollapsed}
            className="app-sider"
            width={240}
          >
            <Menu
              mode="inline"
              className="app-menu"
            >
              <Menu.Item key="home" icon={<HomeOutlined />}>
                <Link to="/">首页</Link>
              </Menu.Item>
              <Menu.Item key="courses" icon={<AppstoreOutlined />}>
                <Link to="/courses">所有课程</Link>
              </Menu.Item>
              {userInfo && userInfo.role === 1 && (
                <>
                  <Menu.SubMenu key="admin" icon={<UserOutlined />} title="管理员功能">
                    <Menu.Item key="publish-course">
                      <Link to="/admin/publish">发布课程</Link>
                    </Menu.Item>
                    <Menu.Item key="upload-video">
                      <Link to="/admin/upload">上传视频</Link>
                    </Menu.Item>
                  </Menu.SubMenu>
                </>
              )}
            </Menu>
          </Sider>

          <Content className="app-content">
            <div className="content-wrapper">
              {/* 登录/注册弹窗 */}
              <Modal
                title={isLoginTab ? "用户登录" : "学生注册"}
                open={loginVisible}
                onCancel={() => {
                  setLoginVisible(false);
                  form.resetFields();
                  registerForm.resetFields();
                }}
                footer={null}
                className="auth-modal"
                width={400}
              >
                <div className="auth-tabs">
                  <Button
                    type={isLoginTab ? "primary" : "text"}
                    onClick={() => setIsLoginTab(true)}
                    className="auth-tab-btn"
                  >
                    登录
                  </Button>
                  <Button
                    type={!isLoginTab ? "primary" : "text"}
                    onClick={() => setIsLoginTab(false)}
                    className="auth-tab-btn"
                  >
                    注册
                  </Button>
                </div>

                {isLoginTab ? (
                  <Form
                    form={form}
                    layout="vertical"
                    initialValues={{ username: 'student789', password: '123456' }}
                    onFinish={handleLogin}
                    className="auth-form"
                  >
                    <Form.Item
                      label="用户名"
                      name="username"
                      rules={[{ required: true, message: '请输入用户名' }]}
                    >
                      <Input size="large" />
                    </Form.Item>
                    <Form.Item
                      label="密码"
                      name="password"
                      rules={[{ required: true, message: '请输入密码' }]}
                    >
                      <Input.Password size="large" />
                    </Form.Item>
                    <Form.Item>
                      <Button type="primary" htmlType="submit" size="large" block>
                        登录
                      </Button>
                    </Form.Item>
                  </Form>
                ) : (
                  <Form
                    form={registerForm}
                    layout="vertical"
                    onFinish={handleRegister}
                    className="auth-form"
                  >
                    <Form.Item
                      label="用户名"
                      name="username"
                      rules={[
                        { required: true, message: '请输入用户名' },
                        { min: 3, message: '用户名长度不能少于3位' }
                      ]}
                    >
                      <Input size="large" placeholder="请输入用户名（3-50位）" />
                    </Form.Item>
                    <Form.Item
                      label="密码"
                      name="password"
                      rules={[
                        { required: true, message: '请输入密码' },
                        { min: 6, message: '密码长度不能少于6位' }
                      ]}
                    >
                      <Input.Password size="large" placeholder="请输入密码（不少于6位）" />
                    </Form.Item>
                    <Form.Item>
                      <Button type="primary" htmlType="submit" size="large" block>
                        注册
                      </Button>
                    </Form.Item>
                  </Form>
                )}
              </Modal>

              {/* 管理员功能区 */}
              {userInfo && userInfo.role === 1 && (
                <>
                  <Card className="admin-card" title="课程管理" extra={<VideoCameraOutlined />}>
                    <Row gutter={[24, 24]}>
                      <Col span={12}>
                        <Card title="发布新课程" className="admin-subcard">
                          <Form
                            form={courseForm}
                            layout="vertical"
                            onFinish={handlePublishCourse}
                            initialValues={{ price: 99.00 }}
                          >
                            <Form.Item
                              label="课程标题"
                              name="title"
                              rules={[{ required: true, message: '请输入课程标题' }]}
                            >
                              <Input placeholder="例如：Go语言从入门到精通" />
                            </Form.Item>
                            <Form.Item
                              label="课程价格"
                              name="price"
                              rules={[{ required: true, message: '请输入课程价格' }]}
                            >
                              <InputNumber min={0.01} precision={2} style={{ width: '100%' }} placeholder="请设置课程价格" />
                            </Form.Item>
                            <Form.Item
                              label="课程简介"
                              name="intro"
                              rules={[{ required: true, message: '请输入课程简介' }]}
                            >
                              <Input.TextArea rows={3} placeholder="请简要描述课程内容" />
                            </Form.Item>
                            <Form.Item label="课程封面">
                              <div className="cover-upload-area">
                                <Upload
                                  beforeUpload={() => false}
                                  onChange={handleCoverChange}
                                  showUploadList={false}
                                  listType="picture-card"
                                  accept="image/*"
                                >
                                  {coverPreview ? (
                                    <div className="cover-preview">
                                      <img src={coverPreview} alt="封面预览" />
                                    </div>
                                  ) : (
                                    <div className="cover-upload-btn">
                                      <UploadOutlined style={{ fontSize: 24, color: '#1890ff' }} />
                                      <div style={{ marginTop: 8 }}>选择封面</div>
                                    </div>
                                  )}
                                </Upload>
                                <Button 
                                  type="primary" 
                                  onClick={handleUploadCover} 
                                  style={{ marginTop: 10 }}
                                  disabled={!coverFile}
                                  block
                                >
                                  上传封面到OSS
                                </Button>
                              </div>
                            </Form.Item>
                            <Form.Item name="cover_oss_key" hidden>
                              <Input />
                            </Form.Item>
                            <Form.Item>
                              <Button type="primary" htmlType="submit" block>
                                发布课程
                              </Button>
                            </Form.Item>
                          </Form>
                        </Card>
                      </Col>
                      <Col span={12}>
                        <Card title="上传课程视频" className="admin-subcard">
                          <Form
                            form={videoForm}
                            layout="vertical"
                            onFinish={handleUploadVideo}
                          >
                            <Form.Item
                              label="关联课程"
                              name="course_id"
                              rules={[{ required: true, message: '请选择要绑定的课程' }]}
                            >
                              <Select placeholder="请选择已发布的课程">
                                {courses.map(course => (
                                  <Option key={course.id} value={course.id}>
                                    {course.title}（¥{course.price.toFixed(2)}）
                                  </Option>
                                ))}
                              </Select>
                            </Form.Item>

                            <Form.Item
                              label="视频标题"
                              name="video_title"
                              rules={[{ required: true, message: '请输入视频标题' }]}
                            >
                              <Input placeholder="例如：第1章 - Go语言环境搭建" />
                            </Form.Item>

                            <Form.Item label="选择视频">
                              <Upload
                                beforeUpload={() => false}
                                onChange={handleVideoChange}
                                showUploadList={false}
                                listType="file"
                                accept="video/*"
                              >
                                <Button icon={<UploadOutlined />} block>
                                  选择视频文件（MP4/AVI等）
                                </Button>
                              </Upload>
                              {videoFile && (
                                <div className="video-file-info">
                                  <p>已选择：{videoFile.name}</p>
                                  <p>大小：{(videoFile.size / 1024 / 1024).toFixed(2)} MB</p>
                                </div>
                              )}
                            </Form.Item>

                            <Form.Item>
                              <Button 
                                type="primary" 
                                htmlType="submit" 
                                block
                                icon={<VideoCameraOutlined />}
                                disabled={!videoFile}
                              >
                                上传视频到OSS
                              </Button>
                            </Form.Item>
                          </Form>
                        </Card>
                      </Col>
                    </Row>
                  </Card>
                </>
              )}

              {/* 课程列表 */}
              <Card className="courses-card" title="课程中心" extra={`共 ${courses.length} 门课程`}>
                {loading ? (
                  <div className="loading-container">
                    <Spin size="large" />
                    <p>正在加载课程数据...</p>
                  </div>
                ) : courses.length > 0 ? (
                  <Row gutter={[24, 24]}>
                    {courses.map(course => (
                      <Col key={course.id} xs={24} sm={12} md={8} lg={6}>
                        <Card
                          hoverable
                          className="course-card"
                          cover={
                            <div className="course-cover">
                              <img 
                                src={course.coverSignedUrl || 'https://via.placeholder.com/300x180?text=Course+Cover'} 
                                alt={course.title}
                              />
                              <div className="course-price">
                                ¥{course.price.toFixed(2)}
                              </div>
                            </div>
                          }
                        >
                          <div className="course-content">
                            <h3 className="course-title">{course.title}</h3>
                            <p className="course-intro">{course.intro}</p>
                            <div className="course-meta">
                              <span className="course-time">发布时间：{course.createdAt || '未知'}</span>
                            </div>
                            <div className="course-actions">
                              <Button
                                type="primary"
                                icon={<ShoppingCartOutlined />}
                                onClick={() => handleBuyCourse(course.id)}
                                block
                                disabled={userInfo?.role === 1}
                                className="buy-btn"
                              >
                                {userInfo?.role === 1 ? '管理员无需购买' : '购买课程'}
                              </Button>
                              <Button
                                icon={<PlayCircleOutlined />}
                                onClick={() => handlePlayCourse(course.id)}
                                block
                                className="play-btn"
                              >
                                播放课程
                              </Button>
                            </div>
                          </div>
                        </Card>
                      </Col>
                    ))}
                  </Row>
                ) : (
                  <div className="empty-courses">
                    <AppstoreOutlined style={{ fontSize: 48, color: '#999', marginBottom: 16 }} />
                    <p>暂无课程数据，请管理员发布课程</p>
                  </div>
                )}
              </Card>

              {/* 视频播放弹窗 */}
              <Modal
                title={`${courses.find(c => c.id === currentCourseId)?.title || '课程视频播放'}`}
                open={videoModalVisible}
                onCancel={() => {
                  setVideoModalVisible(false);
                  setCourseVideos([]);
                  setCurrentPlayUrl('');
                  setCurrentVideoTitle('');
                  setPlayerLoading(false);
                  setCurrentCourseId(null);
                }}
                width="90%"
                footer={[
                  <Button
                    key="close"
                    type="primary"
                    onClick={() => {
                      setVideoModalVisible(false);
                      setCourseVideos([]);
                      setCurrentPlayUrl('');
                      setCurrentVideoTitle('');
                      setPlayerLoading(false);
                      setCurrentCourseId(null);
                    }}
                    className="close-btn"
                  >
                    关闭播放
                  </Button>
                ]}
                destroyOnClose={true}
                maskClosable={false}
                className="video-modal"
              >
                {renderVideoModalContent()}
              </Modal>
            </div>
          </Content>
        </Layout>

        <Footer className="app-footer">
          <div className="footer-content">
            <div className="footer-links">
              <a href="#">关于我们</a>
              <a href="#">帮助中心</a>
              <a href="#">联系我们</a>
              <a href="#">服务条款</a>
              <a href="#">隐私政策</a>
            </div>
            <p className="footer-copyright">
              在线视频学习平台 © 2026 | 基于 React + Go + gRPC + OSS 开发
            </p>
          </div>
        </Footer>
      </Layout>
    </BrowserRouter>
  );
};

export default App;