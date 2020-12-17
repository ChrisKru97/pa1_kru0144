// use crossbeam::sync::WaitGroup;
use rand::Rng;
// use std::env::args;
use std::fmt::{Display, Formatter, Result};
use std::thread::sleep;
use std::time::Duration;

struct Field {
    data: Vec<Vec<bool>>,
    width: usize,
    height: usize,
    next: Vec<Vec<bool>>,
}

impl Field {
    fn new(width: usize, height: usize) -> Field {
        Field {
            width,
            height,
            data: vec![vec![false; width]; height],
            next: vec![vec![false; width]; height],
        }
    }

    fn set(&mut self, x: usize, y: usize, value: bool) {
        self.data[y][x] = value;
    }

    fn set_next(&mut self, x: usize, y: usize, value: bool) {
        self.next[y][x] = value;
    }

    fn prepare_switch(&mut self) {
        for i in 0..self.height {
            for j in 0..self.width {
                self.next[i][j] = self.data[i][j];
            }
        }
    }

    fn switch_to_next(&mut self) {
        let temp = (&self.next).to_vec();
        self.next = (&self.data).to_vec();
        self.data = temp;
    }

    fn is_alive(&self, x: usize, y: usize) -> bool {
        let mut alive: u8 = 0;
        let x_min: i16 = if x > 0 { -1 } else { 0 };
        let x_max: i16 = if x < (self.width - 1) { 2 } else { 1 };
        let y_min: i16 = if y > 0 { -1 } else { 0 };
        let y_max: i16 = if y < (self.height - 1) { 2 } else { 1 };
        for i in x_min..x_max {
            for j in y_min..y_max {
                if !(i == 0 && j == 0)
                    && self.data[(y as i16 + j) as usize][(x as i16 + i) as usize]
                {
                    alive += 1;
                }
            }
        }
        alive == 3 || (self.data[y][x] && alive == 2)
    }
}

impl Display for Field {
    fn fmt(&self, f: &mut Formatter) -> Result {
        let mut s = String::new();
        for row in &self.data {
            let char_vec: Vec<char> = row
                .into_iter()
                .map(|&v| if v { 'O' } else { ' ' })
                .collect();
            let row_string: String = char_vec.into_iter().collect();
            s.push_str(&row_string);
        }
        write!(f, "{}", s)
    }
}

struct Life {
    width: usize,
    height: usize,
    field: Field,
    // thread_count: usize,
}

impl Life {
    fn new(width: usize, height: usize) -> Life {
        let mut field = Field::new(width, height);
        let mut rng = rand::thread_rng();
        let mut x: usize;
        let mut y: usize;
        for _ in 0..(width * height / 5) {
            x = rng.gen_range(0, width);
            y = rng.gen_range(0, height);
            field.set(x, y, true);
        }
        Life {
            width,
            height,
            field,
            // thread_count,
        }
    }

    fn update(&mut self, index: usize) {
        for i in (index..self.height).step_by(1) {
            for j in 0..self.width {
                self.field.set_next(j, i, self.field.is_alive(j, i));
            }
        }
    }

    fn start(&mut self) {
        loop {
            // let wg = WaitGroup::new();
            println!("{}", self.field);
            self.field.prepare_switch();
            self.update(0);
            // for i in 0..self.thread_count {
            //     let wg = wg.clone();
            //     spawn(move || {
            //         self.update(i);
            //         drop(wg);
            //     });
            // }
            // wg.wait();
            self.field.switch_to_next();
            sleep(Duration::from_millis(66));
        }
    }
}

fn main() {
    // let mut thread_count: usize = 8;
    // if let Some(x) = args().into_iter().nth(1) {
    //     thread_count = x.parse::<usize>().unwrap()
    // }
    if let Some((w, h)) = term_size::dimensions() {
        let mut life = Life::new(w, h);
        life.start();
    } else {
        println!("Unable to get term size")
    }
}
