import api from './client';

export interface Course {
  id: string;
  title: string;
  description: string;
  language: string;
  level: string;
  price: number;
  thumbnail_url: string;
}

export const courseService = {
  getAllCourses: async (params?: any) => {
    const response = await api.get('/api/v1/courses', { params });
    return response.data;
  },

  getCourse: async (id: string) => {
    const response = await api.get(`/api/v1/courses/${id}`);
    return response.data;
  },

  enrollInCourse: async (id: string) => {
    const response = await api.post(`/api/v1/courses/${id}/enroll`);
    return response.data;
  },

  getMyCourses: async () => {
    const response = await api.get('/api/v1/courses/my-courses');
    return response.data;
  },
};

export const authService = {
  login: async (email: string, password: string) => {
    const response = await api.post('/api/auth/login', { email, password });
    return response.data;
  },

  register: async (email: string, password: string, role: string) => {
    const response = await api.post('/api/auth/register', {
      email,
      password,
      confirmPassword: password,
      role,
    });
    return response.data;
  },
};
