**<h1 align="center">PRACTICA 2</h1>**

**<h3 align="center">Control y monitoreo de procesos, fase 2</h3>**

### Datos de los estudiantes

<div align="center">

| Nombre                      | Carné     |
| --------------------------- | --------- |
| Adrian Samuel Molina Cabrera    | 201903850 |
| Carlos Raúl Campos Meléndez | 201800639 |
| Kevin Josué Calderón Peraza | 201902714 |

</div>


# MANUAL TECNICO
* Descripción de la configuración de la Máquina virtual
<p align="center">
  <img src="https://github.com/kevcalderon/SO1_201902714/blob/master/Practica2/img/vms.png" width="600">
</p>

* Instrucciones de instalación de Arch Linux
- pacstrap -K /mnt base linux linux-firmware
- genfstab -U /mnt >> /mnt/etc/fstab
- arch-chroot /mnt
- ln -sf /usr/share/zoneinfo/Región/Ciudad /etc/localtime
- hwclock --systohc
- locale-gen
- /etc/locale.conf
- LANG=es_ES.UTF-8
- passwd

Luego se instala la version grafica de gnome para hacer funcionarlo graficamente


* Instrucciones de instalación de los paquetes necesarios para el correcto funcionamiento de la aplicación.

sudo pacman -S wget
wget https://go.dev/dl/go1.20.4.linux-amd64.tar.gz
sudo tar -xvf go1.20.4.linux-amd64.tar.gz -C /usr/local/
export PATH=$PATH:/usr/local/go/bin


* Descripción de los cambios que fueron necesarios para el correcto funcionamiento de los módulos en la nueva distrubucion.


# MANUAL DE USUARIO
