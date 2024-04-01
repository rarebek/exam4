User uchun tokenlar 2 xil bo'ldi:
  1) byAdmin - admin o'zi user create qilsa, tokenning 'subject'i byAdmin ga o'zgaradi.
  2) byRegister - user o'zi ro'yxatdan o'tsa, tokenning 'subject'i byRegister ga o'zgaradi

Userning GetAll methodi "https://localhost:5555/v1/users{page}/{limit}?page=1&limit=10" barcha userlarni, yozgan postlari va o'sha postdagi  barcha commentlarni pagination bilan olib beradi.


Userdagi 'Create' va 'Register' metodlarida, access tokenni responseda qaytarmasdan, headerga avtomatik Authorizatsiya qilib qo'ydim.

